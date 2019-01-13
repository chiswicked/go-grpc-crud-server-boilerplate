package item

import (
	"testing"

	"github.com/chiswicked/go-grpc-crud-server-boilerplate/model"
	pb "github.com/chiswicked/go-grpc-crud-server-boilerplate/protobuf"
)

func TestModelToPb(t *testing.T) {
	if modelToPb(nil) != nil {
		t.Fail()
	}
	item := modelToPb(&model.Item{ID: "abc-123", Name: "Test Name"})
	if item.Id != "abc-123" || item.Name != "Test Name" {
		t.Fail()
	}
}

func TestPbToModel(t *testing.T) {
	if pbToModel(nil) != nil {
		t.Fail()
	}
	item := pbToModel(&pb.Item{Id: "abc-123", Name: "Test Name"})
	if item.ID != "abc-123" || item.Name != "Test Name" {
		t.Fail()
	}
}
