package repository

func Init(url string) error {
	r, err := New(url)
	if err != nil {
		return err
	}

	defer r.Close()

	return r.Write()
}
