package email_verification

import (
	"context"
	"net/http"
	"os"

	"github.com/ScottRMackenzie/Go-Lang-Publishing-Platform/db"
)

func ValidateToken(ctx context.Context, token string) (string, error) {
	var userID string
	err := db.Pool.QueryRow(ctx, "SELECT user_id FROM email_verifications WHERE verification_token = $1", token).Scan(&userID)
	if err != nil {
		return "", err
	}

	return userID, nil
}

func DeleteToken(ctx context.Context, token string) error {
	_, err := db.Pool.Exec(ctx, "DELETE FROM email_verifications WHERE verification_token = $1", token)
	if err != nil {
		return err
	}

	return nil
}

func SetUserVerified(ctx context.Context, userID string) error {
	_, err := db.Pool.Exec(ctx, "UPDATE users SET is_verified = true WHERE id = $1", userID)
	if err != nil {
		return err
	}

	return nil
}

func VerifyEmailHandler(w http.ResponseWriter, r *http.Request) {
	token := r.PathValue("token")
	userID, err := ValidateToken(context.Background(), token)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusBadRequest)
	}

	err = SetUserVerified(context.Background(), userID)
	if err != nil {
		http.Error(w, "Failed to verify email", http.StatusInternalServerError)
	}

	err = DeleteToken(context.Background(), token)
	if err != nil {
		http.Error(w, "Failed to delete token", http.StatusInternalServerError)
	}

	success_url := os.Getenv("BASE_URL") + "/create-account/confirm-email/success"
	http.Redirect(w, r, success_url, http.StatusSeeOther)
}
