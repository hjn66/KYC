function checkLogin(nonce) {
  $.get("/logins?nounce=" + nonce, function(data) {
    var login = JSON.parse(data);
    if (login.Nounce) {
      // alert("Data Loaded: " + data);
      $("#ticket").addClass("hidden");
      $("#qr").addClass("hidden");
      $("#nounce").addClass("hidden");

      $("#FirstName").removeClass("hidden");
      if (login.CheckFirstName) {
        $("#FirstName").addClass("green");
      } else {
        $("#FirstName").addClass("red");
      }
      $("#FirstName").text("First Name : " + login.FirstName);

      $("#LastName").removeClass("hidden");
      if (login.CheckLastName) {
        $("#LastName").addClass("green");
      } else {
        $("#LastName").addClass("red");
      }
      $("#LastName").text("Last Name : " + login.LastName);

      $("#Image").removeClass("hidden");
      if (login.CheckLastName) {
        $("#Image").addClass("green");
      } else {
        $("#Image").addClass("red");
      }
      $("#Image").attr("src", "data:image/png;base64," + login.Image);
    } else {
      setTimeout(checkLogin, 1000, nonce);
    }
  });
}
