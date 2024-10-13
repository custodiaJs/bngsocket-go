package bngsocket

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
)

// Accept wartet auf eingehende Channel-Anfragen und registriert eine neue Channel-Sitzung.
func (o *BngConnChannelListener) Accept() (*BngConnChannel, error) {
	// Überprüfen, ob der Acceptor offen ist.
	if !o.waitOfAccepting.IsOpen() {
		return nil, fmt.Errorf("BngConnChannelListener->Accept[0]: accepting not possible")
	}

	// Auf neue Acceptor-Anfragen warten.
	acceptorRequest, ok := o.waitOfAccepting.Read()
	if !ok {
		return nil, fmt.Errorf("BngConnChannelListener->Accept[1]: cant read from chan")
	}

	// Eine neue eindeutige ID für die Channel-Sitzung erzeugen.
	id := strings.ReplaceAll(uuid.New().String(), "-", "")

	// Eine neue Channel-Sitzung registrieren.
	channlObject, err := o.socket._RegisterNewChannelSession(id)
	if err != nil {
		return nil, fmt.Errorf("BngConnChannelListener->BngConnChannel: %s", err.Error())
	}

	// Die Antwort an den anfragenden Channel zurücksenden.
	if err := responseNewChannelSession(o.socket, acceptorRequest.requestChannelid, id); err != nil {
		return nil, fmt.Errorf("BngConnChannelListener->Accept: %s", err.Error())
	}

	// Das registrierte Channel-Objekt zurückgeben.
	return channlObject, nil
}

// processIncommingSessionRequest verarbeitet eingehende Anfragen zur Eröffnung einer neuen Channel-Sitzung.
func (o *BngConnChannelListener) processIncommingSessionRequest(requestChannelId string, requestedChannelid string) error {
	// Mutex für den Zugriffsschutz auf das Objekt verwenden.
	o.mu.Lock()
	defer o.mu.Unlock()

	// Überprüfen, ob der Kanal offen ist.
	if !o.waitOfAccepting.IsOpen() {
		return fmt.Errorf("BngConnChannelListener->processIncommingSessionRequest[0]: channel ist closed")
	}

	// Ein neues Request-Objekt für die Channel-Anfrage erstellen.
	reqObj := &bngConnAcceptingRequest{
		requestedChannelId: requestedChannelid, // ID des angeforderten Channels
		requestChannelid:   requestChannelId,   // ID des anfragenden Channels
	}

	// Das Request-Objekt in den Kanal schreiben.
	if !o.waitOfAccepting.Enter(reqObj) {
		return fmt.Errorf("BngConnChannelListener->processIncommingSessionRequest[1]: chan was closed")
	}

	// Keine Fehler aufgetreten, Rückgabe nil.
	return nil
}

// Close schließt den Channel Listener (derzeit ohne Implementierung).
func (o *BngConnChannelListener) Close() error {
	// Placeholder für die Schließlogik, derzeit keine Operation.
	return nil
}