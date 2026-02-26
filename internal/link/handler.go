package link

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/JBK2116/phakelinks/types"
	"github.com/gorilla/mux"
)

// LinkConn holds the database connection for link-related queries.
type LinkConn struct {
	logger *slog.Logger
}

// NewLinkConn() creates a new LinkConn with the provided database connection.
func NewLinkConn(logger *slog.Logger) *LinkConn {
	return &LinkConn{
		logger: logger,
	}
}

// RegisterRoutes() registers all routes for the LinkConn struct
func (linkConn *LinkConn) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/links", linkConn.handleCreateLink).Methods("POST")
}

// handleCreateLink() handles the business logic for creating a new link
func (linkConn *LinkConn) handleCreateLink(writer http.ResponseWriter, request *http.Request) {
	var dto types.CreateLinkDTO

	if err := json.NewDecoder(request.Body).Decode(&dto); err != nil {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(writer).Encode(map[string]string{"error": err.Error()})
		linkConn.logger.Error("Error decoding CreateLinkDTO payload", slog.Any("error", err))
		return
	}
	defer request.Body.Close()

	if errStruct := ValidateCreateLinkDTO(dto); errStruct != nil {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(writer).Encode(errStruct)
		linkConn.logger.Info("Invalid CreateLinkDTO payload", slog.Any("error", errStruct))
		return
	}

	var returnDTO types.ReturnLinkDTO
	if dto.Mode == string(types.Educational) {
		randPhishingTechnique := GetRandomPhishingTechnique(dto.Exclude)
		explanationDTO, err := GetEducationalAISummary(randPhishingTechnique, dto.Link)
		if err != nil {
			writer.Header().Set("Content-Type", "application/json")
			writer.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(writer).Encode(map[string]string{"error": err.Error()})
			linkConn.logger.Info("Error creating explanationDTO", slog.Any("error", err.Error()))
			return
		}
		returnDTO.FakeLink = explanationDTO.FakeLink
		returnDTO.Technique = explanationDTO.Technique
		returnDTO.Explanation = explanationDTO.Explanation
	} else {
		prankDTO, err := GetPrankLink(dto.Link)
		if err != nil {
			writer.Header().Set("Content-Type", "application/json")
			writer.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(writer).Encode(map[string]string{"error": err.Error()})
			linkConn.logger.Info("Error creating prankDTO", slog.Any("error", err.Error()))
			return
		}
		returnDTO.FakeLink = prankDTO.Link
	}
	returnDTO.Link = dto.Link
	returnDTO.Mode = dto.Mode
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(returnDTO)
}
