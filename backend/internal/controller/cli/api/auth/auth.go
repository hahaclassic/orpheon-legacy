package auth_cli_ctrl

import (
	"bufio"
	"context"
	"fmt"
	"os"

	"github.com/hahaclassic/orpheon/backend/internal/controller/cli/session"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	usecase "github.com/hahaclassic/orpheon/backend/internal/domain/usecases/auth"
	"github.com/hahaclassic/orpheon/backend/pkg/cmdrouter"
)

type AuthController struct {
	authService usecase.AuthService
}

func NewAuthController(authService usecase.AuthService) *AuthController {
	authController := &AuthController{authService: authService}

	return authController
}

func (c *AuthController) Menu() []cmdrouter.OptionHandler {
	return []cmdrouter.OptionHandler{
		{
			Name: "Login",
			Run:  c.login,
		},
		{
			Name: "Register",
			Run:  c.register,
		},
		{
			Name: "Logout",
			Run:  c.logout,
		},
		{
			Name: "Update Password",
			Run:  c.updatePassword,
		},
		{
			Name: "Refresh Token",
			Run:  c.RefreshToken,
		},
	}
}

func (c *AuthController) login(ctx context.Context) error {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Enter login: ")
	scanner.Scan()
	login := scanner.Text()

	fmt.Print("Enter password: ")
	scanner.Scan()
	password := scanner.Text()

	credentials := &entity.UserCredentials{
		Login:    login,
		Password: password,
	}

	tokens, err := c.authService.Login(ctx, credentials)
	if err != nil {
		return fmt.Errorf("login failed: %w", err)
	}

	claims, err := c.authService.GetClaims(ctx, tokens.Access)
	if err != nil {
		return fmt.Errorf("failed to get claims: %w", err)
	}

	session.StartSession(claims, tokens)

	fmt.Printf("Login successful. You are logged in as '%s'\n", credentials.Login)

	return nil
}

func (c *AuthController) register(ctx context.Context) error {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Enter login: ")
	scanner.Scan()
	login := scanner.Text()

	fmt.Print("Enter password: ")
	scanner.Scan()
	password := scanner.Text()

	credentials := &entity.UserCredentials{
		Login:    login,
		Password: password,
	}

	tokens, err := c.authService.RegisterUser(ctx, credentials)
	if err != nil {
		return fmt.Errorf("register failed: %w", err)
	}

	claims, err := c.authService.GetClaims(ctx, tokens.Access)
	if err != nil {
		return fmt.Errorf("failed to get claims: %w", err)
	}

	session.StartSession(claims, tokens)

	fmt.Printf("Register successful. You are logged in as '%s'\n", credentials.Login)

	return nil
}

func (c *AuthController) updatePassword(ctx context.Context) error {
	if !session.IsAuthenticated() {
		fmt.Println("You are not logged in. Please login first.")
		return nil
	}

	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Enter old password: ")
	scanner.Scan()
	oldPassword := scanner.Text()

	fmt.Print("Enter new password: ")
	scanner.Scan()
	newPassword := scanner.Text()

	fmt.Print("Confirm new password: ")
	scanner.Scan()
	newPasswordConfirm := scanner.Text()

	if newPassword != newPasswordConfirm {
		return fmt.Errorf("new password and confirm new password do not match")
	}

	passwords := &entity.UserPasswords{
		Old: oldPassword,
		New: newPassword,
	}

	claims := session.Claims()

	err := c.authService.UpdatePassword(ctx, claims.UserID, passwords)
	if err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	fmt.Println("Password updated successfully")

	return nil
}

func (c *AuthController) logout(ctx context.Context) error {
	err := c.authService.Logout(ctx, session.Tokens().Refresh)
	if err != nil {
		return fmt.Errorf("failed to logout: %w", err)
	}

	session.EndSession()

	fmt.Println("Logged out successfully")

	return nil
}

func (c *AuthController) RefreshToken(ctx context.Context) error {
	tokens, err := c.authService.RefreshTokens(ctx, session.Tokens().Refresh)
	if err != nil {
		return fmt.Errorf("failed to refresh token: %w", err)
	}

	claims, err := c.authService.GetClaims(ctx, tokens.Access)
	if err != nil {
		return fmt.Errorf("failed to get claims: %w", err)
	}

	session.StartSession(claims, tokens)

	fmt.Println("Tokens refreshed successfully")

	return nil
}
