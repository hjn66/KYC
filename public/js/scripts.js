function checkLogin(nonce) {
  $.get("/logins?nonce=" + nonce, function(data) {
    var login = JSON.parse(data);
    if (login.Nonce) {
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
    console.log(register);
    if (register.Nonce) {
      // alert("Data Loaded: " + data);
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
