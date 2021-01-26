package container

import (
	"context"
)

func CheckImageSource(ctx context.Context, instance interface{}) bool {
	return true
}

func CopyImageToBackUpRegistry(ctx context.Context, instance interface{}) (string, error){
	return "", nil
}

func UpdateImageFromResource(ctx context.Context, instance interface{}, image string) error {
	return nil
}

func ValidateRedeployedApplication(ctx context.Context, instance interface{}, magfe string) error {
	return nil
}
