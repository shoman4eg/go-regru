package regru

import (
	"context"
	"testing"
)

func TestFolderService(t *testing.T) {
	ctx := context.Background()
	content := `{"folder_content":[{"domain_name":"test1.ru","service_id":"1000","servtype":"domain"},{"domain_name":"test2.ru","service_id":"1001","servtype":"domain"}]}`
	services := `{"folder_content":[{"domain_name":"test1.ru","service_id":"1000"}],"services":[{"dname":"test01.ru","result":"success","service_id":"123456","servtype":"domain"},{"dname":"test11.ru","error_code":"DOMAIN_NOT_FOUND","error_params":{"domain_name":"test11.ru"},"result":"Domain test11.ru not found or not owned by You"}]}`
	byName := FolderServicesOptions{FolderIdent: FolderIdent{FolderName: "test_folder_name"}, Services: []ServiceIdent{{DomainName: "test1.ru"}, {DomainName: "test2.ru"}}, ReturnFolderContents: true}
	byNameIn := map[string]any{"folder_name": "test_folder_name", "services": []map[string]string{{"domain_name": "test1.ru"}, {"domain_name": "test2.ru"}}, "return_folder_contents": 1}
	hasContent := func(v any) bool {
		a := v.(*FolderServicesAnswer)
		return len(a.FolderContent) > 0 && a.FolderContent[0].ServiceID == 1000
	}
	hasServices := func(v any) bool {
		a := v.(*FolderServicesAnswer)
		return hasContent(v) && a.Services[0].ServiceID == 123456 && ErrorCode(a.Services[1].Err()) == "DOMAIN_NOT_FOUND"
	}
	tests := []struct {
		path, answer string
		want         any
		call         func(*Client) (any, error)
		ok           func(any) bool
	}{
		{"/folder/nop", `{"id":"-1","name":"test_folder_name"}`, map[string]any{"folder_id": 123456, "folder_name": "test_folder_name"},
			func(c *Client) (any, error) {
				r, _, err := c.Folder.Nop(ctx, &FolderIdent{FolderID: 123456, FolderName: "test_folder_name"})
				return r, err
			},
			func(v any) bool { a := v.(*FolderAnswer); return a.ID == -1 && a.Name == "test_folder_name" }},
		{"/folder/create", `{}`, map[string]any{"folder_name": "test_folder_name"},
			func(c *Client) (any, error) {
				r, _, err := c.Folder.Create(ctx, &FolderIdent{FolderName: "test_folder_name"})
				return r, err
			}, nil},
		{"/folder/remove", `{}`, map[string]any{"folder_id": 123456},
			func(c *Client) (any, error) {
				r, _, err := c.Folder.Remove(ctx, &FolderIdent{FolderID: 123456})
				return r, err
			}, nil},
		{"/folder/rename", content, map[string]any{"folder_name": "test_folder_name", "new_folder_name": "new_test_folder_name"},
			func(c *Client) (any, error) {
				r, _, err := c.Folder.Rename(ctx, &FolderRenameOptions{FolderIdent: FolderIdent{FolderName: "test_folder_name"}, NewFolderName: "new_test_folder_name"})
				return r, err
			}, hasContent},
		{"/folder/get_services", content, map[string]any{"folder_id": 12345},
			func(c *Client) (any, error) {
				r, _, err := c.Folder.GetServices(ctx, &FolderIdent{FolderID: 12345})
				return r, err
			}, hasContent},
		{"/folder/add_services", services, byNameIn,
			func(c *Client) (any, error) { r, _, err := c.Folder.AddServices(ctx, &byName); return r, err }, hasServices},
		{"/folder/remove_services", services, byNameIn,
			func(c *Client) (any, error) { r, _, err := c.Folder.RemoveServices(ctx, &byName); return r, err }, hasServices},
		{"/folder/replace_services", services, byNameIn,
			func(c *Client) (any, error) { r, _, err := c.Folder.ReplaceServices(ctx, &byName); return r, err }, hasServices},
		{"/folder/move_services", services,
			map[string]any{"folder_id": 12345, "new_folder_id": 1234567, "services": []map[string]int{{"service_id": 1000}, {"service_id": 1001}}},
			func(c *Client) (any, error) {
				r, _, err := c.Folder.MoveServices(ctx, &FolderMoveServicesOptions{
					FolderServicesOptions: FolderServicesOptions{FolderIdent: FolderIdent{FolderID: 12345}, Services: []ServiceIdent{{ServiceID: 1000}, {ServiceID: 1001}}},
					NewFolderID:           1234567,
				})
				return r, err
			}, hasServices},
	}
	c, mux := setup(t)
	for _, tt := range tests {
		handle(t, mux, tt.path, tt.want, tt.answer)
		if got, err := tt.call(c); err != nil {
			t.Errorf("%s: %v", tt.path, err)
		} else if tt.ok != nil && !tt.ok(got) {
			t.Errorf("%s: got %+v", tt.path, got)
		}
	}
}
