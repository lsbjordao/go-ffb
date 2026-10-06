package ffb

import "context"

func (c *Client) FetchFamilies(ctx context.Context) ([]string, error) {
	body, err := c.get(ctx, "families")
	if err != nil {
		return nil, err
	}
	return decodeStrings(body)
}

func (c *Client) FetchGeneraOfFamily(ctx context.Context, family string) ([]string, error) {
	body, err := c.get(ctx, "genus", family)
	if err != nil {
		return nil, err
	}
	return decodeStrings(body)
}
