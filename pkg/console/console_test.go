package console

import "testing"

func Test_fmtURL(t *testing.T) {
	type args struct {
		token        string
		targetRegion string
	}

	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "use default region",
			args: args{
				token: "Fizz",
			},
			want:    "https://signin.aws.amazon.com/federation?Action=login&Destination=https%3A%2F%2Feu-west-1.console.aws.amazon.com%2F&Issuer=aws-console&SigninToken=Fizz",
			wantErr: false,
		},
		{
			name: "overwrite default region London region",
			args: args{
				token:        "buzz",
				targetRegion: "eu-west-2",
			},
			want:    "https://signin.aws.amazon.com/federation?Action=login&Destination=https%3A%2F%2Feu-west-2.console.aws.amazon.com%2F&Issuer=aws-console&SigninToken=buzz",
			wantErr: false,
		},
		{
			name: "let's try a CN region",
			args: args{
				token:        "helloCN,",
				targetRegion: "cn-north-1",
			},
			want:    "https://signin.amazonaws.cn/federation?Action=login&Destination=https%3A%2F%2Fcn-north-1.console.amazonaws.cn%2F&Issuer=aws-console&SigninToken=helloCN%2C",
			wantErr: false,
		},
		{
			name: "fail on missing token",
			args: args{
				targetRegion: "us-west-1",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := fmtURL(tt.args.token, tt.args.targetRegion)
			if (err != nil) != tt.wantErr {
				t.Errorf("fmtURL() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if got != tt.want {
				t.Errorf("fmtURL() got = %v, want %v", got, tt.want)
			}
		})
	}
}
