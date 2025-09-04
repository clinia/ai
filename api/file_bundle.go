package api

// BundleType disambiguates single output vs multi-file archive.
type BundleType int

const (
	BundleTypeUnknown BundleType = iota
	BundleTypeSingleFile
	BundleTypeArchive
)
