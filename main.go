package main

import (
	"fmt"
	"time"
)

//Region UseCase
var isDoorOpenedInput IsDoorOpenedCondition
var isDoorClosedInput IsDoorClosedCondition
var userOpen UserInput
var userClose UserInput

func main() {
	isDoorOpenedInput = true
	isDoorClosedInput = false
	userOpen = false
	userClose = false

	stateOpened := State{Name: "opened"}
	stateClosed := State{Name: "closed"}
	stateOpening := State{Name: "opening", Action: openDoor}
	stateClosing := State{Name: "closing", Action: closeDoor}
	states := []State{stateOpened, stateClosed, stateOpening, stateClosing}

	eventOpen := Event{Name: "open", Condition: &userOpen}
	eventClose := Event{Name: "close", Condition: &userClose}
	eventSensorOpened := Event{Name: "sensor opened", Condition: &isDoorOpenedInput}
	eventSensorClosed := Event{Name: "sensor closed", Condition: &isDoorClosedInput}
	events := []Event{eventOpen, eventClose, eventSensorOpened, eventSensorClosed}

	transitions := []Transition{
		Transition{CurrentState: stateOpened, Event: eventClose, TargetState: stateClosing},
		Transition{CurrentState: stateClosing, Event: eventOpen, TargetState: stateClosing},
		Transition{CurrentState: stateClosing, Event: eventSensorClosed, TargetState: stateClosed},
		Transition{CurrentState: stateClosed, Event: eventOpen, TargetState: stateOpening},
		Transition{CurrentState: stateOpening, Event: eventClose, TargetState: stateClosing},
		Transition{CurrentState: stateOpening, Event: eventSensorOpened, TargetState: stateOpened},
	}

	fmt.Println("States: ", states)
	fmt.Println("Events: ", events)
	fmt.Println("Transitions: ", transitions)

	fmt.Println("Running state machine")
	doorSM := NewStateMachine(transitions)
	go doorSM.Run()

	time.Sleep(3 * time.Second)
	userClose = true

	time.Sleep(1 * time.Second)
	userClose = false

	time.Sleep(3 * time.Second)
	userOpen = true
	time.Sleep(1 * time.Second)
	userOpen = false

	for {
	}
}

//Actions
func openDoor() {
	fmt.Println("Opening Door")
	time.Sleep(time.Second)
	isDoorOpenedInput = true
	isDoorClosedInput = false
}

func closeDoor() {
	fmt.Println("Closing Door")
	time.Sleep(time.Second)
	isDoorOpenedInput = false
	isDoorClosedInput = true
}

//Condition handling
type UserInput bool

func (in UserInput) Test() bool {
	return bool(in)
}

type IsDoorOpenedCondition bool

func (c IsDoorOpenedCondition) Test() bool {
	return bool(isDoorOpenedInput) && !bool(userClose)
}

type IsDoorClosedCondition bool

func (c IsDoorClosedCondition) Test() bool {
	return bool(isDoorClosedInput) && !bool(userOpen)
}

//Region State-Machine
type State struct {
	Name   string
	Action func()
}

type ConditionInterface interface {
	Test() bool
}

type Event struct {
	Name      string
	Condition ConditionInterface
}

type Transition struct {
	CurrentState State
	Event        Event
	TargetState  State
}

type StateMachine struct {
	transitions  []Transition
	currentState State
}

func NewStateMachine(transitions []Transition) StateMachine {
	return StateMachine{transitions, transitions[0].CurrentState}
}

func (sm *StateMachine) Run() error {
	for {
		fmt.Printf("Current State %s\n", sm.currentState)
		eventsForState := sm.getEventsForCurrentState()

		var triggeredEvent *Event
		for _, event := range eventsForState {
			fmt.Println("Testing ", event.Name)
			if event.Condition.Test() {
				if triggeredEvent != nil {
					panic("At least two conditions hit for same state")
				}
				triggeredEvent = &event
			}
		}

		if triggeredEvent != nil {
			fmt.Println("Triggered Event", *triggeredEvent)
			sm.executeTransition(*triggeredEvent)
			fmt.Printf("Now moved into %s state\n", sm.currentState)
		}

		if sm.currentState.Action != nil {
			sm.currentState.Action()
		}
		time.Sleep(time.Second)
	}
	return nil
}

func (sm StateMachine) getEventsForCurrentState() []Event {
	events := []Event{}
	for _, t := range sm.transitions {
		if t.CurrentState.Name == sm.currentState.Name {
			events = append(events, t.Event)
		}
	}
	return events
}

func (sm StateMachine) getTransition(currState State, event Event) Transition {
	for _, t := range sm.transitions {
		if t.CurrentState.Name == currState.Name && t.Event.Name == event.Name {
			return t
		}
	}
	panic("Could not find transition")
}

func (sm *StateMachine) executeTransition(event Event) {
	t := sm.getTransition(sm.currentState, event)
	sm.currentState = t.TargetState
}
