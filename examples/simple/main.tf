resource "gmail_label" "test" {
  name                    = "my-test"
  label_list_visibility   = "labelShow"
  message_list_visibility = "show"
}
