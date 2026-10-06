package ffb

import "context"

func (c *Client) FetchSpeciesOfGenus(ctx context.Context, genus string) ([]string, error) {
	body, err := c.get(ctx, "species", "genus", genus)
	if err != nil {
		return nil, err
	}
	return decodeStrings(body)
}
