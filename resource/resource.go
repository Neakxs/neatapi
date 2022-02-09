package resource

type Field interface {
	GetTag(s string) string
	GetName() string
}

type Resource interface {
	GetPublicNames(f Field) []string
	GetPrivateName(f Field) string
}

type FieldValidater interface {
	ValidateField(f Field) error
}
