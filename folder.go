package regru

import "context"

// FolderService handles functions of the "folder" category.
type FolderService service

// FolderAnswer is returned by Nop, Create and Remove.
type FolderAnswer struct {
	ID   FlexInt `json:"id"`
	Name string  `json:"name"`
}

// Nop checks that a folder exists (folder/nop).
func (s *FolderService) Nop(ctx context.Context, opts *FolderIdent) (*FolderAnswer, *Response, error) {
	return call[FolderAnswer](ctx, s.client, "folder/nop", opts)
}

// Create creates a folder named opts.FolderName (folder/create).
func (s *FolderService) Create(ctx context.Context, opts *FolderIdent) (*FolderAnswer, *Response, error) {
	return call[FolderAnswer](ctx, s.client, "folder/create", opts)
}

// Remove deletes a folder (folder/remove).
func (s *FolderService) Remove(ctx context.Context, opts *FolderIdent) (*FolderAnswer, *Response, error) {
	return call[FolderAnswer](ctx, s.client, "folder/remove", opts)
}

// FolderContentItem is a service in a folder.
type FolderContentItem struct {
	DomainName string  `json:"domain_name"`
	ServiceID  FlexInt `json:"service_id"`
	ServType   string  `json:"servtype"`
}

// FolderServicesAnswer is returned by Rename and the *Services functions.
// FolderContent is set by Rename, GetServices and when ReturnFolderContents is requested.
type FolderServicesAnswer struct {
	FolderContent []FolderContentItem `json:"folder_content"`
	Services      []ServiceItem       `json:"services"`
}

// FolderRenameOptions are the parameters of folder/rename.
type FolderRenameOptions struct {
	FolderIdent
	NewFolderName string `json:"new_folder_name"`
}

// Rename renames a folder (folder/rename).
func (s *FolderService) Rename(ctx context.Context, opts *FolderRenameOptions) (*FolderServicesAnswer, *Response, error) {
	return call[FolderServicesAnswer](ctx, s.client, "folder/rename", opts)
}

// GetServices lists the services in a folder (folder/get_services).
func (s *FolderService) GetServices(ctx context.Context, opts *FolderIdent) (*FolderServicesAnswer, *Response, error) {
	return call[FolderServicesAnswer](ctx, s.client, "folder/get_services", opts)
}

// FolderServicesOptions are the parameters of folder/add_services, remove_services and replace_services.
type FolderServicesOptions struct {
	FolderIdent
	Services             []ServiceIdent `json:"services"`
	ReturnFolderContents Bool           `json:"return_folder_contents,omitempty"`
}

// AddServices adds services to a folder (folder/add_services).
func (s *FolderService) AddServices(ctx context.Context, opts *FolderServicesOptions) (*FolderServicesAnswer, *Response, error) {
	return call[FolderServicesAnswer](ctx, s.client, "folder/add_services", opts)
}

// RemoveServices removes services from a folder (folder/remove_services).
func (s *FolderService) RemoveServices(ctx context.Context, opts *FolderServicesOptions) (*FolderServicesAnswer, *Response, error) {
	return call[FolderServicesAnswer](ctx, s.client, "folder/remove_services", opts)
}

// ReplaceServices replaces all services in a folder with the given ones (folder/replace_services).
func (s *FolderService) ReplaceServices(ctx context.Context, opts *FolderServicesOptions) (*FolderServicesAnswer, *Response, error) {
	return call[FolderServicesAnswer](ctx, s.client, "folder/replace_services", opts)
}

// FolderMoveServicesOptions are the parameters of folder/move_services.
type FolderMoveServicesOptions struct {
	FolderServicesOptions
	NewFolderID   int64  `json:"new_folder_id,omitempty"`
	NewFolderName string `json:"new_folder_name,omitempty"`
}

// MoveServices moves services from one folder to another (folder/move_services).
func (s *FolderService) MoveServices(ctx context.Context, opts *FolderMoveServicesOptions) (*FolderServicesAnswer, *Response, error) {
	return call[FolderServicesAnswer](ctx, s.client, "folder/move_services", opts)
}
