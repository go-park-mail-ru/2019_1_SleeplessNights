package message

import "fmt"

func (m *Message) IsValid() bool {
	logger.Info("Logger Entered IsValid ")
	fmt.Println(m.Payload)
	switch m.Title {
	case Ready:
		{
			return true
		}
	case GoTo:
		{

			st, ok := m.Payload.(map[string]interface{})
			if !ok {
				logger.Error("Message validator, Title=GO_TO, error:interface->Answer casting error")
				return false
			}
			if _, ok := st["x"]; !ok {
				return false
			}
			if _, ok := st["y"]; !ok {
				return false
			}
			return true
		}
	case ClientAnswer:
		{
			st, ok := m.Payload.(map[string]interface{})
			if !ok {
				logger.Error("Message validator, Title=ClientAnswer, error:interface->Answer casting error")
				return false
			}
			if _, ok := st["answer_id"]; !ok {
				return false
			}

			return true
		}
	case Leave:
		{
			return true
		}

	case Continue:
		{
			return true
		}
	case ChangeOpponent:
		{
			return true
		}
	case Quit:
		{
			return true
		}
	case State:
		{
			return true
		}
	case NotDesiredPack:
		{
			return true
		}
	default:
		return false
	}
}
