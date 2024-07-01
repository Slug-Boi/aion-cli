package form

import (
	"context"
	"fmt"

	"google.golang.org/api/forms/v1"
)

// CreateService creates a new form service
// https://pkg.go.dev/google.golang.org/api/forms/v1#hdr-Creating_a_client
func CreateService() (*forms.Service, error) {
  
  // TODO: Implement authentication
  ctx := context.Background()
  formsService, err := forms.NewService(ctx)
  if err != nil {
    return nil, fmt.Errorf("unable to create forms service: %v", err)
  } 

  return formsService, nil
}

// GetForm retrieves a form and returns it as marsheled JSON
func GetForm(service *forms.Service, formID string) ([]byte, error) {

  ListRespHeader, err := service.Forms.Responses.List(formID).Do()
  if err != nil {
    return nil, fmt.Errorf("unable to retrieve form responses: %v", err)
  }

  return ListRespHeader.MarshalJSON()
}