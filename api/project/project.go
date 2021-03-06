package project

type Project struct {
	Id                   string            `json:"id" bson:"_id,omitempty"`
	Name                 string            `json:"name"`
	Banner               string            `json:"banner" bson:"banner"`
	Description          string            `json:"description"`
	Tags                 []string          `json:"tags"`
	Owner                string            `json:"owner"`
	ImageProvider        string            `json:"imageProvider" bson:"imageprovider"`
	ImageProviderOptions map[string]string `json:"imageProviderOptions" bson:"imageprovideroptions"`
}
