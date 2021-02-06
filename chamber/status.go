package chamber

type Status struct {
    Id           string `json:"id"`
    Status       string `json:"status"`
    UptimeMillis int    `json:"uptimeMillis"`
    Port         int    `json:"port"`

    Spec Spec `json:"spec"`
}
