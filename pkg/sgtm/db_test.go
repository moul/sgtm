package sgtm

import (
	"testing"

	"github.com/stretchr/testify/require"
	"moul.io/sgtm/pkg/sgtmpb"
)

func TestDBUserCreate(t *testing.T) {
	db := TestingDB(t)

	tests := []struct {
		name            string
		input           *sgtmpb.User
		expectedOutput  *sgtmpb.User
		shouldHaveError bool
	}{
		{"nil", nil, nil, true},
		/*{
			"empty",
			&sgtmpb.User{},
			nil,
			true,
		},*/
		{
			"manfred",
			&sgtmpb.User{Firstname: "Manfred", Lastname: "Touron", Email: "m@42.am", Username: "moul"},
			&sgtmpb.User{Firstname: "Manfred", Lastname: "Touron", Email: "m@42.am", Username: "moul"},
			false,
		},
		{
			"manfred2",
			&sgtmpb.User{Firstname: "Manfred2", Lastname: "Touron2", Email: "m@42.am2", Username: "moul2"},
			&sgtmpb.User{Firstname: "Manfred2", Lastname: "Touron2", Email: "m@42.am2", Username: "moul2"},
			false,
		},
		{
			"manfred:again",
			&sgtmpb.User{Firstname: "Manfred", Lastname: "Touron", Email: "m@42.am", Username: "moul"},
			nil,
			true,
		},
		{
			"manfred2:again",
			&sgtmpb.User{Firstname: "Manfred2", Lastname: "Touron2", Email: "m@42.am2", Username: "moul2"},
			nil,
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := db.Create(tt.input).Error
			if tt.shouldHaveError {
				require.Error(t, err)
				require.Nil(t, tt.expectedOutput)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, tt.input)
			require.NotZero(t, tt.input.ID)
			require.NotZero(t, tt.input.CreatedAt)
			require.NotZero(t, tt.input.UpdatedAt)
			require.NotNil(t, tt.expectedOutput)

			// copy dynamic fields before comparison
			tt.expectedOutput.ID = tt.input.ID
			tt.expectedOutput.CreatedAt = tt.input.CreatedAt
			tt.expectedOutput.UpdatedAt = tt.input.UpdatedAt
			require.Equal(t, tt.input, tt.expectedOutput)

			//fmt.Println(godev.PrettyJSON(tt.input))
		})
	}
}