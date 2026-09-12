package constant

type UploadPurpose string

const (
	UploadPurposeAvatar       UploadPurpose = "avatar"
	UploadPurposeLogo         UploadPurpose = "logo"
	UploadPurposeService      UploadPurpose = "service"
	UploadPurposePaymentProof UploadPurpose = "payment-proof"
)

var MimetypeWhitelist = map[string]bool{
	"image/jpg":  true,
	"image/jpeg": true,
	"image/png":  true,
	"image/webp": true,
}

func GetAssetType(mimetype string) string {
	switch mimetype {
	case "image/jpg", "image/jpeg", "image/png", "image/webp":
		return "image"
	default:
		return "auto"
	}
}
