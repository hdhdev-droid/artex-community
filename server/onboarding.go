package server

import "net/http"

// Only return the metadata used by the checklist. The full MCP configuration
// may contain credentials in environment variables, arguments or endpoint URLs.
// A separate route also fails closed against an older backend.
func (s *Server) onboardingMCP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store")
	pg := s.pg(w)
	if pg == nil {
		return
	}
	servers, err := pg.ListMCP()
	if err != nil {
		writeErr(w, 503, "onboarding_status_unavailable")
		return
	}
	items := make([]map[string]any, 0, len(servers))
	for _, server := range servers {
		items = append(items, map[string]any{"id": server.ID, "name": server.Name, "enabled": server.Enabled})
	}
	writeJSON(w, 200, map[string]any{"servers": items})
}
