package prometheus

import "github.com/prometheus/client_golang/prometheus"

var (
	CreateMovieApiPingCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "create_movie_api_ping_counter",
		Help: "Number of pings made to the Create Movie API endpoint",
	})
	CreateActorApiPingCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "create_actor_api_ping_counter",
		Help: "Number of pings made to the Create Actor API endpoint",
	})
	RegisterApiPingCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "register_api_ping_counter",
		Help: "Number of pings made to the Register API endpoint",
	})
	LoginApiPingCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "login_api_ping_counter",
		Help: "Number of pings made to the Login API endpoint",
	})
	UpdateMovieApiPingCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "update_movie_api_ping_counter",
		Help: "Number of pings made to the Update Movie API endpoint",
	})
	UpdateActorApiPingCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "update_actor_api_ping_counter",
		Help: "Number of pings made to the Update Actor API endpoint",
	})
	PatchActorApiPingCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "patch_actor_api_ping_counter",
		Help: "Number of pings made to the Patch Actor API endpoint",
	})
	PatchMovieApiPingCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "patch_movie_api_ping_counter",
		Help: "Number of pings made to the Patch Movie API endpoint",
	})
	DeleteMovieApiPingCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "delete_movie_api_ping_counter",
		Help: "Number of pings made to the Delete Movie API endpoint",
	})
	DeleteActorApiPingCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "delete_actor_api_ping_counter",
		Help: "Number of pings made to the Delete Actor API endpoint",
	})
	ReadAllActorApiPingCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "read_all_actor_api_ping_counter",
		Help: "Number of pings made to the Read All Actor API endpoint",
	})
	ReadOneActorApiPingCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "read_one_actor_api_ping_counter",
		Help: "Number of pings made to the Read One Actor API endpoint",
	})
	ReadOneMovieApiPingCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "read_one_movie_api_ping_counter",
		Help: "Number of pings made to the Read One Movie API endpoint",
	})
	ReadAllMovieApiPingCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "read_all_movie_api_ping_counter",
		Help: "Number of pings made to the Read All Movie API endpoint",
	})
)

func Init() {
	prometheus.MustRegister(CreateMovieApiPingCounter)
	prometheus.MustRegister(CreateActorApiPingCounter)
	prometheus.MustRegister(RegisterApiPingCounter)
	prometheus.MustRegister(LoginApiPingCounter)
	prometheus.MustRegister(UpdateMovieApiPingCounter)
	prometheus.MustRegister(UpdateActorApiPingCounter)
	prometheus.MustRegister(PatchActorApiPingCounter)
	prometheus.MustRegister(PatchMovieApiPingCounter)
	prometheus.MustRegister(DeleteMovieApiPingCounter)
	prometheus.MustRegister(DeleteActorApiPingCounter)
	prometheus.MustRegister(ReadAllActorApiPingCounter)
	prometheus.MustRegister(ReadOneActorApiPingCounter)
	prometheus.MustRegister(ReadOneMovieApiPingCounter)
	prometheus.MustRegister(ReadAllMovieApiPingCounter)
}
