function checkLogin(nonce) {
  $.get("/logins?nonce=" + nonce, function(data) {
    var login = JSON.parse(data);
    if (login.nonce) {
      // alert("Data Loaded: " + data);
      $("#ticket").addClass("hidden");
      $("#qr").addClass("hidden");
      $("#nonce").addClass("hidden");

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

function checkRegister(nonce) {
  $.get("/registers?nonce=" + nonce, function(data) {
    var register = JSON.parse(data);
    if (register.nonce) {
      // alert("Data Loaded: " + data);
      $("#ticket").addClass("hidden");
      $("#qr").addClass("hidden");
      $("#nonce").addClass("hidden");
      $("#registerForm").removeClass("hidden");

      $("#FirstName").removeClass("hidden");
      $("#FirstName").val(register.User.FirstName);

      $("#LastName").removeClass("hidden");
      $("#LastName").val(register.User.LastName);
      $("#BirthDate").val(register.User.BirthDate);

      $("#Image").removeClass("hidden");
      $("#Image").attr("src", "data:image/png;base64," + register.User.Photo);
    } else {
      setTimeout(checkRegister, 1000, nonce);
    }
  });
}
