package image

import (
	"errors"

	"github.com/containers/image/types"
)

// memoryImage is a mostly-implementation of types.Image assembled from data
// created in memory, used primarily as a return value of types.Image.UpdatedImage
// as a way to carry various structured information in a type-safe and easy-to-use way.
// Note that this _only_ carries the immediate metadata; it is _not_ a stand-alone
// collection of all related information, e.g. there is no way to get layer blobs
// from a memoryImage.
type memoryImage struct {
	genericManifest
	serializedManifest []byte // A private cache for Manifest()
}

func memoryImageFromManifest(m genericManifest) types.Image {
	return &memoryImage{
		genericManifest:    m,
		serializedManifest: nil,
	}
}

// Reference returns the reference used to set up this source, _as specified by the user_
// (not as the image itself, or its underlying storage, claims).  This can be used e.g. to determine which public keys are trusted for this image.
func (i *memoryImage) Reference() types.ImageReference {
	// It would really be inappropriate to return the ImageReference of the image this was based on.
	return nil
}

// Close removes resources associated with an initialized UnparsedImage, if any.
func (i *memoryImage) Close() {
}

// Manifest is like ImageSource.GetManifest, but the result is cached; it is OK to call this however often you need.
func (i *memoryImage) Manifest() ([]byte, string, error) {
	if i.serializedManifest == nil {
		m, err := i.genericManifest.serialize()
		if err != nil {
			return nil, "", err
		}
		i.serializedManifest = m
	}
	return i.serializedManifest, i.genericManifest.manifestMIMEType(), nil
}

// Signatures is like ImageSource.GetSignatures, but the result is cached; it is OK to call this however often you need.
func (i *memoryImage) Signatures() ([][]byte, error) {
	// Modifying an image invalidates signatures; a caller asking the updated image for signatures
	// is probably confused.
	return nil, errors.New("Internal error: Image.Signatures() is not supported for images modified in memory")
}

// Inspect returns various information for (skopeo inspect) parsed from the manifest and configuration.
func (i *memoryImage) Inspect() (*types.ImageInspectInfo, error) {
	return inspectManifest(i.genericManifest)
}

// IsMultiImage returns true if the image's manifest is a list of images, false otherwise.
func (i *memoryImage) IsMultiImage() bool {
	return false
}
