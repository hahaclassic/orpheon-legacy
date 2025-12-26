package license_cli_ctrl

import (
	"bufio"
	"context"
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/controller/cli/output"
	"github.com/hahaclassic/orpheon/backend/internal/controller/cli/session"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	usecase "github.com/hahaclassic/orpheon/backend/internal/domain/usecases/content/license"
	"github.com/hahaclassic/orpheon/backend/pkg/cmdrouter"
)

type LicenseController struct {
	licenseService usecase.LicenseService
}

func NewLicenseController(licenseService usecase.LicenseService) *LicenseController {
	return &LicenseController{
		licenseService: licenseService,
	}
}

func (c *LicenseController) Menu() []cmdrouter.OptionHandler {
	return []cmdrouter.OptionHandler{
		{
			Name: "Create",
			Run:  c.createLicense,
		},
		{
			Name: "Get by ID",
			Run:  c.getLicenseByID,
		},
		{
			Name: "List licenses",
			Run:  c.getAllLicenses,
		},
		{
			Name: "Update",
			Run:  c.updateLicense,
		},
		{
			Name: "Delete",
			Run:  c.deleteLicense,
		},
	}
}

func (c *LicenseController) createLicense(ctx context.Context) error {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Enter license name: ")
	scanner.Scan()
	name := scanner.Text()

	fmt.Print("Enter license description: ")
	scanner.Scan()
	description := scanner.Text()

	license := &entity.License{
		Title:       name,
		Description: description,
	}

	claims := session.Claims()

	err := c.licenseService.CreateLicense(ctx, claims, license)
	if err != nil {
		return fmt.Errorf("failed to create license: %w", err)
	}

	fmt.Println("License created successfully")
	return nil
}

func (c *LicenseController) getLicenseByID(ctx context.Context) error {
	var id string

	fmt.Print("Enter license ID: ")
	if _, err := fmt.Scan(&id); err != nil {
		return fmt.Errorf("failed to read license ID: %w", err)
	}

	licenseID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("failed to parse license ID: %w", err)
	}

	license, err := c.licenseService.GetLicenseByID(ctx, licenseID)
	if err != nil {
		return fmt.Errorf("failed to get license: %w", err)
	}

	output.PrintLicense(license)
	return nil
}

func (c *LicenseController) getAllLicenses(ctx context.Context) error {
	licenses, err := c.licenseService.GetAllLicenses(ctx)
	if err != nil {
		return fmt.Errorf("failed to list licenses: %w", err)
	}

	output.PrintLicenses(licenses)

	return nil
}

func (c *LicenseController) updateLicense(ctx context.Context) error {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Enter license ID: ")
	scanner.Scan()
	id := scanner.Text()

	fmt.Print("Enter new license title: ")
	scanner.Scan()
	title := scanner.Text()

	fmt.Print("Enter new license description: ")
	scanner.Scan()
	description := scanner.Text()

	licenseID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("failed to parse license ID: %w", err)
	}

	license := &entity.License{
		ID:          licenseID,
		Title:       title,
		Description: description,
	}

	claims := session.Claims()

	err = c.licenseService.UpdateLicense(ctx, claims, license)
	if err != nil {
		return fmt.Errorf("failed to update license: %w", err)
	}

	fmt.Println("License updated successfully")
	return nil
}

func (c *LicenseController) deleteLicense(ctx context.Context) error {
	var id string

	fmt.Print("Enter license ID: ")
	if _, err := fmt.Scan(&id); err != nil {
		return fmt.Errorf("failed to read license ID: %w", err)
	}

	licenseID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("failed to parse license ID: %w", err)
	}

	claims := session.Claims()

	err = c.licenseService.DeleteLicense(ctx, claims, licenseID)
	if err != nil {
		return fmt.Errorf("failed to delete license: %w", err)
	}

	fmt.Println("License deleted successfully")
	return nil
}
