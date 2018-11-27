function checkLogin(nonce) {
  $.get("/logins?nonce=" + nonce, function(data) {
    var login = JSON.parse(data);
    if (login.Nonce) {
      $("#qr").addClass("hidden");
      $("#TicketB").addClass("hidden");
      $("#TicketD").addClass("hidden");
      $("#NonceB").addClass("hidden");
      $("#NonceD").addClass("hidden");

      $("#LoginForm").removeClass("hidden");

      $("#GUID").text("GUID: " + login.GUID);
      $("#LoginDate").text("Login Date: " + login.LoginDate);
      $("#FirstName").text("First Name: " + login.FirstName);
      $("#LastName").text("Last Name: " + login.LastName);
      $("#Image").attr("src", "data:image/png;base64," + login.Image);
      if (login.CheckFirstName) {
        $("#FirstName").addClass("green");
      } else {
        $("#FirstName").addClass("red");
      }
      if (login.CheckLastName) {
        $("#LastName").addClass("green");
      } else {
        $("#LastName").addClass("red");
      }
      if (login.CheckImage) {
        $(".badge").attr("src", "images/true.png");
      } else {
        $(".badge").attr("src", "images/false.png");
      }
    } else {
      setTimeout(checkLogin, 1000, nonce);
    }
  });
}

function checkRegister(nonce) {
  $.get("/registers?nonce=" + nonce, function(data) {
    var register = JSON.parse(data);
    console.log(register);
    if (register.Nonce) {
      $("#TicketB").addClass("hidden");
      $("#TicketD").addClass("hidden");
      $("#NonceB").addClass("hidden");
      $("#NonceD").addClass("hidden");
      $("#qr").addClass("hidden");
      $("#registerForm").removeClass("hidden");

      $("#NationalIDForm").val(register.User.NationalID);
      $("#FirstNameForm").val(register.User.FirstName);
      $("#LastNameForm").val(register.User.LastName);
      $("#BirthDateForm").val(register.User.BirthDate);
      $("#PublicKeyForm").val(register.User.PublicKey);
      $("#PhotoForm").val(register.User.Photo);
      $("#Nonce").val(register.Nonce);

      $("#RegisterDate").text("Register Date: " + register.RegisterDate);
      $("#NationalID").text("National ID: " + register.User.NationalID);
      $("#FirstName").text("First Name: " + register.User.FirstName);
      $("#LastName").text("Last Name: " + register.User.LastName);
      $("#BirthDate").text("Birth Date: " + register.User.BirthDate);
      $("#PublicKey").text(register.User.PublicKey);
      $("#Image").attr("src", "data:image/png;base64," + register.User.Photo);
    } else {
      setTimeout(checkRegister, 1000, nonce);
    }
  });
}

function submitForm(action) {
  $("#Action").val(action);
  $("#registerForm").submit();
}
