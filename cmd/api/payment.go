package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/checkout/session"
	"github.com/stripe/stripe-go/v82/webhook"
)

func (app *application) createCheckoutSession(w http.ResponseWriter, r *http.Request) {
	var paymentData struct {
		Price      float64 `json:"price"`
		Currency   string  `json:"currency"`
		Email      string  `json:"email"`
		MenteeName string  `json:"menteeName"`
		MentorName string  `json:"mentorName"`
		MentorID   int     `json:"mentorId"`
	}

	if err := json.NewDecoder(r.Body).Decode(&paymentData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		log.Printf("json.Decode: %v", err)
		return
	}

	stripe.Key = "sk_test_51NpubpSF2Iapc3CYMensOJgnQjo4anfwi9MNLOFIjNkOBYRxzEP8gMctadHwISPfAERy31iKNejTs50cRCu1bCxV00NycUfZ06"

	domain := "http://localhost:5173/profile"
	priceInRupees := int64(paymentData.Price * 100)

	params := &stripe.CheckoutSessionParams{
		CustomerEmail:            stripe.String(paymentData.Email),
		SubmitType:               stripe.String("book"),
		BillingAddressCollection: stripe.String("auto"),
		ShippingAddressCollection: &stripe.CheckoutSessionShippingAddressCollectionParams{
			AllowedCountries: stripe.StringSlice([]string{
				"IN",
			}),
		},
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			&stripe.CheckoutSessionLineItemParams{
				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
					Currency: stripe.String("inr"),
					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
						Name: stripe.String("Mentorship Session with " + paymentData.MentorName),
					},
					UnitAmount: stripe.Int64(priceInRupees),
				},
				Quantity: stripe.Int64(1),
			},
		},
		Mode:       stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL: stripe.String(domain + "?success=true"),
		CancelURL:  stripe.String(domain + "?canceled=true"),
		Metadata: map[string]string{
			"mentorId": strconv.Itoa(paymentData.MentorID),
		},
	}

	s, err := session.New(params)
	if err != nil {
		http.Error(w, "Failed to create checkout session", http.StatusInternalServerError)
		log.Printf("session.New: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"url": s.URL,
	})
}

func (app *application) handleWebhook(w http.ResponseWriter, r *http.Request) {
	const MaxBodyBytes = int64(65536)
	r.Body = http.MaxBytesReader(w, r.Body, MaxBodyBytes)
	payload, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		log.Printf("io.ReadAll: %v", err)
		return
	}

	// Verify webhook signature
	event, err := webhook.ConstructEvent(payload, r.Header.Get("Stripe-Signature"), app.config.stripeWebhook)
	if err != nil {
		http.Error(w, "Webhook signature verification failed", http.StatusBadRequest)
		log.Printf("webhook.ConstructEvent: %v", err)
		return
	}

	// Handle the checkout.session.completed event
	if event.Type == "checkout.session.completed" {
		var checkoutSession stripe.CheckoutSession
		err := json.Unmarshal(event.Data.Raw, &checkoutSession)
		if err != nil {
			http.Error(w, "Error parsing webhook JSON", http.StatusInternalServerError)
			log.Printf("json.Unmarshal: %v", err)
			return
		}

		// Extract mentorId from metadata
		mentorID, ok := checkoutSession.Metadata["mentorId"]
		if !ok {
			log.Printf("mentorId not found in metadata")
			http.Error(w, "mentorId not found", http.StatusInternalServerError)
			return
		}

		// Call the API to mark the meeting as paid using PUT method
		apiURL := "http://localhost:8080/v1/meetings/paid/" + mentorID
		req, err := http.NewRequest(http.MethodPut, apiURL, nil)
		if err != nil {
			log.Printf("Failed to create PUT request: %v", err)
			http.Error(w, "Failed to update payment status", http.StatusInternalServerError)
			return
		}

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Printf("Failed to call paid API: %v", err)
			http.Error(w, "Failed to update payment status", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			log.Printf("Paid API returned non-OK status: %v", resp.Status)
			http.Error(w, "Failed to update payment status", http.StatusInternalServerError)
			return
		}

		log.Printf("Successfully marked meeting as paid for mentorId: %s", mentorID)
	}

	w.WriteHeader(http.StatusOK)
}
