package room

import (
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/TheGame/game_field"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/TheGame/messge"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/TheGame/player"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/TheGame/questions"
	"sync"
)

//Комната - тот объект, который инкапсулирует в себе всю рботу с игровыми механиками
//Здесь нам нужно внутри игрового цикла слать получать сообщения игроков о том, какие ходы они делают,
//прощитывать новые игровые ситуации и отправлять игрокам сообщения об изменении игровой ситуации

//Из неигровых задач комната должна уметь
//* Собираться и пересобираться, не выкидывая игроков, если оба решили сыграть ещё партию вместе или это лобби
//* Поддерживать обработку отвалившегося игрока

type Room struct {
	p1          player.Player
	p2          player.Player
	field       game_field.GameField
	questions   questions.QuestionPack
	IsAvailable bool       //Флаг, того, что комната уже наполнена и балансировщику сюда не надо стучаться
	mu          sync.Mutex //Добавление игрока в комнату - конкурентная операция, поэтому нужен мьютекс
	//Если не знаете, что это такое, то погуглите (для любого языка), об этом написано много, но, обычно, довольно сложно
	//Если по-простому, то это типа стоп-сигнала для всех остальных потоков, который можно включить,
	//сделать всё, что нужно, пока тебе никто не мешает, и выключить обратно
}

func (r *Room)TryJoin(p player.Player)(success bool) {
	//Здесь нам нужно под мьютексом проверить наличие свободных мест. Варианты:
	//1. 2 места свободно -> занимаем первое место
	//2. Свободно 1 место -> занимаем место, поднимаем флаг недоступности комнаты, начинаем игровой процесс
	//TODO develop
	return
}

func (r *Room)BuildEnv() {
	//Процедура должна пересоздавать игровое поле, запрашивать новый список тем из БД и готовить комнату к новой игре
	//При этом она должна уметь работать асинхронно и не выбрасывать пользователей из комнаты во время работы
	//TODO develop
}

func (r *Room)notify(msg messge.Message) {
	//Процедуры должна отправить сообщение всем игрокам комнаты
	//TODO develop
}

func (r *Room)grantGodMod(p player.Player, token []byte)  {
	//РЕАЛИЗОВЫВАТЬ ПОСЛЕДНЕЙ
	//Это чисто техническая процедура, она нужна не для реальных игроков, а, в основном, для ботов, которые должны знать
	//правильный ответ, чтобы отвечать верно более чем на 25% вопросов
	//Принцип работы следующий:
	//1. У игрока есть токен на получение всех овтветов (бот будет запрашивать его у сервера, из другого пакета)
	//2. Игрок отправляет сообщение с запросом на все ответы и своим токеном
	//3. Получив сообщение, комната запускает эту функцию
	//4. Здесь мы проверяем валидность токена, и возвращаем в сообщении игроку матрицу правильных ответов
	//5. ВАЖНО! Конретно это сообщение надо отправлять напрямую конкретному игроку, а не через notify
	//TODO develop
}

func (r *Room)startMatch() {
	//Эта процедура запускает игровой процесс
	//Здесь мы будем слушать все сообщения пользователей асинхронно и складывать их в очередь для обработки
	//В цикле мы будем обрабатывать все входные сообщения, выполнять нашу бизнес логику (менять значение таймера,
	//давать пользователям вопросы и т.д.)
	//Все сообщения мы будем складывать в очередь на отправку и отправлять всю очередь каждые 0.5 сек
	//(цифра примерная, может поменяться и должна быть вынесена в костанту)
	//TODO develop
}