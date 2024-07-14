package hexapawn

import (
	"os"
	"testing"
)

func TestSave(t *testing.T) {
	// Test case 1: MachineFile is empty
	machine := &Machine{}
	err := machine.Save()
	if err == nil {
		t.Errorf("Expected error when MachineFile is empty, but got nil")
	}

	// Test case 2: MachineFile is not empty
	machine = &Machine{MachineFile: "test_machine.json"}
	err = machine.Save()
	if err != nil {
		t.Errorf("Unexpected error when MachineFile is not empty: %v", err)
	}
	// Check that the file was created
	if _, err := os.Stat("test_machine.json"); os.IsNotExist(err) {
		t.Errorf("Expected file to be created, but it was not")
	}
	// Clean up the test file
	os.Remove("test_machine.json")
}
