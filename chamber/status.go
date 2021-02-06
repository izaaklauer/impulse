package chamber

type Status struct {
    Id           string `json:"id"`
    Status       string `json:"status"`
    CreatedTimeMillis int64    `json:"createdTimeMillis"`
    Port         int    `json:"port"`

    Spec Spec `json:"spec"`
}
