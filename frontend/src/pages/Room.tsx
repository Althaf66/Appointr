import React, { useState, useEffect, useRef } from 'react';
import { useNavigate, useParams } from 'react-router-dom';

export const Room = () => {
    const { roomId } = useParams<{ roomId: string }>();
    const navigate = useNavigate();
    const [isConnected, setIsConnected] = useState(false);
    const [isInitiator, setIsInitiator] = useState(false);
    const [remotePeerId, setRemotePeerId] = useState<string | null>(null);
    const localVideoRef = useRef<HTMLVideoElement>(null);
    const remoteVideoRef = useRef<HTMLVideoElement>(null);
    const peerConnectionRef = useRef<RTCPeerConnection | null>(null);
    const localStreamRef = useRef<MediaStream | null>(null);
    
    // Initialize the meeting room
    useEffect(() => {
      const initRoom = async () => {
        try {
          // Get user media
          const stream = await navigator.mediaDevices.getUserMedia({
            video: true,
            audio: true,
          });
          
          localStreamRef.current = stream;
          
          // Display local video
          if (localVideoRef.current) {
            localVideoRef.current.srcObject = stream;
          }
          
          // Either create or join a room
          const userId = localStorage.getItem('userId') || 'user-' + Math.random().toString(36).substring(2, 7);
          localStorage.setItem('userId', userId);
  
          // Try to join the room first
          const joinResponse = await fetch(`http://localhost:8080/video/join-room/${roomId}`, {
            method: 'POST',
            headers: {
              'Content-Type': 'application/json',
            },
            body: JSON.stringify({ userId }),
          });
          
          if (joinResponse.ok) {
            const joinData = await joinResponse.json();
            setIsInitiator(false);
            setRemotePeerId(joinData.otherPeerId);
            await setupPeerConnection(userId, joinData.otherPeerId, false);
            setIsConnected(true);
          } else if (joinResponse.status === 404) {
            // Room doesn't exist, create it
            const createResponse = await fetch('http://localhost:8080/video/create-room', {
              method: 'POST',
              headers: {
                'Content-Type': 'application/json',
              },
              body: JSON.stringify({ roomId, userId }),
            });
            
            if (createResponse.ok) {
              setIsInitiator(true);
              setIsConnected(true);
              await setupPeerConnection(userId, null, true);
              
              // Poll for peer join events when you're the initiator
              startPollingForPeers(userId);
            } else {
              const errorData = await createResponse.text();
              throw new Error(`Failed to create room: ${errorData}`);
            }
          } else {
            // Add this block to handle other error status codes
            const errorData = await joinResponse.text();
            throw new Error(`Failed to join room: ${errorData}`);
          }
        } catch (error) {
          console.error('Error initializing room:', error);
          alert('Error connecting to room: ' + (error instanceof Error ? error.message : 'Unknown error'));
        }
      };
  
      initRoom();
  
      // Cleanup
      return () => {
        if (localStreamRef.current) {
          localStreamRef.current.getTracks().forEach(track => track.stop());
        }
        if (peerConnectionRef.current) {
          peerConnectionRef.current.close();
        }
      };
    }, [roomId]);

    // Poll for new peers joining the room
    const startPollingForPeers = (userId: string) => {
      // In a real implementation, you would use WebSockets instead of polling
      const intervalId = setInterval(async () => {
        try {
          const response = await fetch(`http://localhost:8080/video/room-status/${roomId}`, {
            method: 'GET',
            headers: {
              'Content-Type': 'application/json',
            }
          });
          
          if (response.ok) {
            const data = await response.json();
            // If we found a new peer that isn't us and we don't already have a remotePeerId
            const otherPeers = data.connections.filter((id: string) => id !== userId);
            if (otherPeers.length > 0 && !remotePeerId) {
              const newPeerId = otherPeers[0];
              setRemotePeerId(newPeerId);
              await createOffer(userId, newPeerId);
            }
          }
        } catch (error) {
          console.error('Error polling for peers:', error);
        }
      }, 2000);
      
      // Return function to clean up interval
      return () => clearInterval(intervalId);
    };
  
    const setupPeerConnection = async (userId: string, remotePeerId: string | null, isInitiator: boolean) => {
      // Create RTCPeerConnection with STUN server
      const configuration: RTCConfiguration = {
        iceServers: [
          {
            urls: 'stun:stun.l.google.com:19302',
          },
        ],
      };
  
      const peerConnection = new RTCPeerConnection(configuration);
      peerConnectionRef.current = peerConnection;
  
      // Add local tracks to the connection
      if (localStreamRef.current) {
        localStreamRef.current.getTracks().forEach(track => {
          if (localStreamRef.current) {
            peerConnection.addTrack(track, localStreamRef.current);
          }
        });
      }
  
      // Set up event handlers for the connection
      peerConnection.ontrack = (event) => {
        console.log('Got remote track:', event.streams[0]);
        if (remoteVideoRef.current && event.streams && event.streams[0]) {
          // Important fix: Set remote video stream and ensure it plays
          remoteVideoRef.current.srcObject = event.streams[0];
          remoteVideoRef.current.play().catch(e => console.error("Error playing remote video:", e));
        }
      };
  
      peerConnection.onicecandidate = async (event) => {
        if (event.candidate && remotePeerId) {
          // Send ICE candidate to the server
          try {
            await fetch('http://localhost:8080/video/signal', {
              method: 'POST',
              headers: {
                'Content-Type': 'application/json',
              },
              body: JSON.stringify({
                type: 'candidate',
                candidate: event.candidate.toJSON(),
                userId,
                roomId,
                to: remotePeerId,
              }),
            });
          } catch (error) {
            console.error('Error sending ICE candidate:', error);
          }
        }
      };
  
      peerConnection.oniceconnectionstatechange = () => {
        console.log('ICE connection state:', peerConnection.iceConnectionState);
        if (peerConnection.iceConnectionState === 'connected' || 
            peerConnection.iceConnectionState === 'completed') {
          setIsConnected(true);
        }
      };

      // For non-initiators, we need to prepare to receive an offer
      if (!isInitiator && remotePeerId) {
        // Set up polling for signaling data (offers, answers, candidates)
        startPollingForSignals(userId, remotePeerId);
      }

      return peerConnection;
    };

    // Poll for signaling messages (offers, answers, candidates)
    const startPollingForSignals = (userId: string, otherPeerId: string) => {
      const pollInterval = setInterval(async () => {
        try {
          const response = await fetch(`http://localhost:8080/video/pending-signals/${roomId}/${userId}`, {
            method: 'GET',
            headers: {
              'Content-Type': 'application/json',
            }
          });

          if (response.ok) {
            const signals = await response.json();
            
            for (const signal of signals) {
              if (signal.type === 'offer') {
                await handleRemoteOffer(userId, otherPeerId, signal.sdp);
              } else if (signal.type === 'candidate' && signal.candidate) {
                // Handle ice candidates
                if (peerConnectionRef.current) {
                  try {
                    await peerConnectionRef.current.addIceCandidate(new RTCIceCandidate(signal.candidate));
                    console.log('Added ICE candidate');
                  } catch (err) {
                    console.error('Error adding received ICE candidate', err);
                  }
                }
              }
            }
          }
        } catch (error) {
          console.error('Error checking for signals:', error);
        }
      }, 1000);

      return () => clearInterval(pollInterval);
    };

    // Handle incoming offer when joining a room where someone is waiting
    const handleRemoteOffer = async (userId: string, remotePeerId: string, offerSdp: string) => {
      if (!peerConnectionRef.current) return;

      try {
        await peerConnectionRef.current.setRemoteDescription(new RTCSessionDescription({
          type: 'offer',
          sdp: offerSdp
        }));

        const answer = await peerConnectionRef.current.createAnswer();
        await peerConnectionRef.current.setLocalDescription(answer);

        // Send the answer back
        await fetch('http://localhost:8080/video/signal', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({
            type: 'answer',
            sdp: answer.sdp,
            userId,
            roomId,
            to: remotePeerId,
          }),
        });
      } catch (error) {
        console.error('Error handling offer:', error);
      }
    };
  
    const createOffer = async (userId: string, remotePeerId: string) => {
      if (!peerConnectionRef.current) return;
  
      try {
        const offer = await peerConnectionRef.current.createOffer();
        await peerConnectionRef.current.setLocalDescription(offer);
  
        // Send the offer to the server
        await fetch('http://localhost:8080/video/signal', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({
            type: 'offer',
            sdp: offer.sdp,
            userId,
            roomId,
            to: remotePeerId,
          }),
        });
        
        // After sending the offer, start polling for the answer
        const checkForAnswer = async () => {
          try {
            const response = await fetch(`http://localhost:8080/video/pending-signals/${roomId}/${userId}`, {
              method: 'GET',
              headers: {
                'Content-Type': 'application/json',
              }
            });
            
            if (response.ok) {
              const signals = await response.json();
              
              for (const signal of signals) {
                if (signal.type === 'answer' && peerConnectionRef.current) {
                  const answer = new RTCSessionDescription({
                    type: 'answer',
                    sdp: signal.sdp,
                  });
                  await peerConnectionRef.current.setRemoteDescription(answer);
                  console.log('Successfully set remote description (answer)');
                  return true; // Successfully processed answer
                } else if (signal.type === 'candidate' && signal.candidate && peerConnectionRef.current) {
                  try {
                    await peerConnectionRef.current.addIceCandidate(new RTCIceCandidate(signal.candidate));
                    console.log('Added ICE candidate from poll');
                  } catch (err) {
                    console.error('Error adding received ICE candidate', err);
                  }
                }
              }
            }
            return false; // No answer processed
          } catch (error) {
            console.error('Error checking for answer:', error);
            return false;
          }
        };
        
        // Poll for answer a few times
        let attempts = 0;
        const maxAttempts = 10;
        const pollForAnswer = setInterval(async () => {
          attempts++;
          const gotAnswer = await checkForAnswer();
          if (gotAnswer || attempts >= maxAttempts) {
            clearInterval(pollForAnswer);
          }
        }, 1000);
      } catch (error) {
        console.error('Error creating offer:', error);
      }
    };
  
    const leaveRoom = () => {
      // Stop all tracks
      if (localStreamRef.current) {
        localStreamRef.current.getTracks().forEach(track => track.stop());
      }
  
      // Close the peer connection
      if (peerConnectionRef.current) {
        peerConnectionRef.current.close();
      }
  
      // Navigate back to home
      navigate('/meeting');
    };

    const toggleMute = () => {
      if (localStreamRef.current) {
        const audioTracks = localStreamRef.current.getAudioTracks();
        if (audioTracks.length > 0) {
          const isEnabled = audioTracks[0].enabled;
          audioTracks[0].enabled = !isEnabled;
        }
      }
    };

    const toggleVideo = () => {
      if (localStreamRef.current) {
        const videoTracks = localStreamRef.current.getVideoTracks();
        if (videoTracks.length > 0) {
          const isEnabled = videoTracks[0].enabled;
          videoTracks[0].enabled = !isEnabled;
        }
      }
    };
  
    return (
      <div className="flex flex-col items-center min-h-screen bg-gray-100 p-4">
        <div className="w-full max-w-4xl bg-white rounded-lg shadow-md p-4 mb-4">
          <div className="flex justify-between items-center mb-4">
            <h1 className="text-xl font-bold">Room: {roomId}</h1>
            <div>
              <span className="mr-4">
                Status: {isConnected ? 'Connected' : 'Connecting...'}
              </span>
              <button
                onClick={leaveRoom}
                className="bg-red-500 text-white py-1 px-4 rounded-md hover:bg-red-600"
              >
                Leave
              </button>
            </div>
          </div>
  
          <div className="flex flex-col md:flex-row space-y-4 md:space-y-0 md:space-x-4">
            <div className="flex-1">
              <h2 className="text-lg font-semibold mb-2">Your Video</h2>
              <div className="bg-black rounded-lg overflow-hidden aspect-video">
                <video
                  ref={localVideoRef}
                  autoPlay
                  muted
                  playsInline
                  className="w-full h-full object-cover"
                />
              </div>
            </div>
  
            <div className="flex-1">
              <h2 className="text-lg font-semibold mb-2">Remote Video</h2>
              <div className="bg-black rounded-lg overflow-hidden aspect-video relative">
                <video
                  ref={remoteVideoRef}
                  autoPlay
                  playsInline
                  className="w-full h-full object-cover"
                />
                {!remotePeerId && (
                  <div className="absolute inset-0 flex items-center justify-center text-white bg-black bg-opacity-70">
                    Waiting for someone to join...
                  </div>
                )}
              </div>
            </div>
          </div>
        </div>
  
        <div className="w-full max-w-4xl bg-white rounded-lg shadow-md p-4">
          <h2 className="text-lg font-semibold mb-2">Meeting Controls</h2>
          <div className="flex justify-center space-x-4">
            <button 
              onClick={toggleMute}
              className="bg-gray-500 text-white py-2 px-4 rounded-full hover:bg-gray-600"
            >
              Mute
            </button>
            <button 
              onClick={toggleVideo}
              className="bg-gray-500 text-white py-2 px-4 rounded-full hover:bg-gray-600"
            >
              Stop Video
            </button>
            <button 
              className="bg-blue-500 text-white py-2 px-4 rounded-full hover:bg-blue-600"
              onClick={() => alert('Screen sharing not implemented yet')}
            >
              Share Screen
            </button>
          </div>
        </div>

        <div className="w-full max-w-4xl bg-white rounded-lg shadow-md p-4 mt-4">
          <h2 className="text-lg font-semibold mb-2">Connection Status</h2>
          <div className="grid grid-cols-2 gap-4">
            <div>
              <p><strong>Your ID:</strong> {localStorage.getItem('userId')}</p>
              <p><strong>Room ID:</strong> {roomId}</p>
              <p><strong>You are:</strong> {isInitiator ? 'Host' : 'Participant'}</p>
            </div>
            <div>
              <p><strong>Remote Peer:</strong> {remotePeerId || 'None yet'}</p>
              <p><strong>ICE State:</strong> {peerConnectionRef.current?.iceConnectionState || 'Not established'}</p>
              <p><strong>Signaling State:</strong> {peerConnectionRef.current?.signalingState || 'Not established'}</p>
            </div>
          </div>
        </div>
      </div>
    );
};