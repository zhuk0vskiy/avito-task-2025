import io.gatling.core.Predef._
import io.gatling.http.Predef._
import scala.concurrent.duration._

class ServerPerSecondLoadSimulation extends Simulation {

  val httpProtocol = http
    .baseUrl("http://172.20.1.1:8080")
    .acceptHeader("application/json")
    .contentTypeHeader("application/json")

  val authScenario1 = scenario("Authentication 1")
    .exec(
      http("Authenticate")
        .post("/api/auth")
        .body(StringBody("""{"username": "testuser1", "password": "password1"}"""))
        .check(status.is(200))
        .check(jsonPath("$.token").saveAs("jwtToken1"))
    )

  val authScenario2 = scenario("Authentication 2")
    .exec(
      http("Authenticate")
        .post("/api/auth")
        .body(StringBody("""{"username": "testuser2", "password": "password2"}"""))
        .check(status.is(200))
        .check(jsonPath("$.token").saveAs("jwtToken2"))
    )

  val infoScenario = scenario("Get Info")
    .exec(authScenario1)
    .exec(
      http("Get Info with JWT")
        .get("/api/info")
        .header("Authorization", "Bearer ${jwtToken1}")
        .check(status.is(200))
    )
  
  val sendScenario1 = scenario("Send Coins 1")
    .exec(authScenario1)
    .exec(
      http("Send Coins to testuser2")
        .post("/api/sendCoin")
        .body(StringBody("""{"toUser": "testuser2", "amount": 10}"""))
        .header("Authorization", "Bearer ${jwtToken1}")
        .check(status.is(200))
    )
    .exec(authScenario2)
    .exec(
      http("Send Coins to testuser1")
        .post("/api/sendCoin")
        .body(StringBody("""{"toUser": "testuser1", "amount": 10}"""))
        .header("Authorization", "Bearer ${jwtToken2}")
        .check(status.is(200))
    )

  // val sendScenario2 = scenario("Send Coins 2")
  //   .exec(authScenario2)
  //   .exec(
  //     http("Send Coins to testuser1")
  //       .post("/api/sendCoin")
  //       .body(StringBody("""{"toUser": "testuser1", "amount": "1"}"""))
  //       .header("Authorization", "Bearer ${jwtToken2}")
  //       .check(status.is(200))
  //   )
  //   .exec(authScenario1)
  //   .exec(
  //     http("Send Coins to testuser2")
  //       .post("/api/sendCoin")
  //       .body(StringBody("""{"toUser": "testuser2", "amount": "1"}"""))
  //       .header("Authorization", "Bearer ${jwtToken1}")
  //       .check(status.is(200))
  //   )

  setUp(
    // authScenario1.inject(atOnceUsers(1)).protocols(httpProtocol),
    // authScenario2.inject(atOnceUsers(1)).protocols(httpProtocol),
    // infoScenario.inject(atOnceUsers(1)).protocols(httpProtocol),
    // sendScenario1.inject(atOnceUsers(1)).protocols(httpProtocol)
    // sendScenario2.inject(atOnceUsers(1)).protocols(httpProtocol)
    infoScenario.inject(
      rampUsersPerSec(10) to 1000 during (60.seconds)
    ).protocols(httpProtocol),
    sendScenario1.inject(
      rampUsersPerSec(10) to 1000 during (60.seconds)
    ).protocols(httpProtocol),
  //   sendScenario2.inject(
  //     rampUsersPerSec(10) to 1000 during (60.seconds)
  //   ).protocols(httpProtocol)
  // )
  )
}