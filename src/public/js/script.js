$(document).ready(function() {
  // do jQuery
  //
  // send request
  $.ajax({
    url: "/api/loginURL", // url where to submit the request
    type: "GET", // type of action POST || GET
    // dataType: "json", // data type
    contentType: "application/json",
    // data: JSON.stringify(data), // post data || get data
    success: function(result) {
      // you can see the result from the console
      // tab of the developer tools
      $("#loginURL").html(result);
    },
    error: function(xhr, resp, text) {
      console.log(xhr, resp, text);
    }
  })

  $("#dynamicField").submit(function(event) {
    event.preventDefault();
    console.log($("#dynamicFieldName").val(), $("#dynamicFieldVal").val());
    $('<input>').attr({
      // type: 'hidden',
      name: $("#dynamicFieldName").val(),
      value: $("#dynamicFieldVal").val()
    }).appendTo('#createApplicationForm');
  });

  $("#childDataForm").submit(function(event) {
    // prepair data for request
    data = processData.bind(this).call();
    sendRequest("/api/addApplicationChildData", data);
    event.preventDefault();
  });

  $("#addUserForm").submit(function(event) {
    // prepair data for request
    data = processData.bind(this).call();
    sendRequest("/api/insertApplicationIntoUser", data);
    event.preventDefault();
  });

  $("#createApplicationForm").submit(function(event) {
    // prepair data for request
    data = processData.bind(this).call();
    sendRequest("/api/createApplication", data);
    event.preventDefault();
  });

  $("#insertForm").submit(function(event) {
    // prepair data for request
    data = processData.bind(this).call();
    sendRequest("/api/insert", data);
    // sendRequestTest("/api/insert", data);
    event.preventDefault();
  });


  $("#updateForm").submit(function(event) {
    // prepair data for request
    data = processData.bind(this).call();
    sendRequest("/api/update", data);
    event.preventDefault();
  });

  $("#validateForm").submit(function(event) {
    // prepair data for request
    data = processData.bind(this).call();
    sendRequest("/api/validate", data);
    event.preventDefault();
  });

  $("#getChildren").submit(function(event) {
    event.preventDefault();
    // prepair data for request
    data = processData.bind(this).call();
    sendRequest("/api/ReturnAppChildren", data);
  });
  $("#getUsers").submit(function(event) {
    event.preventDefault();
    // prepair data for request
    data = processData.bind(this).call();
    sendRequest("/api/ReturnAppUsers", data);
  });

  $("#appList").change(function() {
    $("span").hide();
    var str = "";
    $("#appList option:selected").each(function() {
      str += $(this).text() + " ";
    })
    $("#" + str).show();
  });

  $("span").hide();
  console.log("hello world");
})


// processData takes form data and processes and readies data for ajax post
function processData() {

  var data = {};
  //Gathering the Data
  //and removing undefined keys(buttons)
  $.each(this.elements, function(i, v) {
    var input = $(v);
    data[input.attr("name")] = input.val();
    delete data["undefined"];
  });
  return data;
}

function sendRequest(url, data) {
  // send request
  $.ajax({
    url: url, // url where to submit the request
    type: "POST", // type of action POST || GET
    // dataType: "json", // data type
    contentType: "application/json",
    data: JSON.stringify(data), // post data || get data
    success: function(result) {
      // you can see the result from the console
      // tab of the developer tools
      alert(result);
      console.log(result);
      $(".result").html(result);
    },
    error: function(xhr, resp, text) {
      console.log(xhr, resp, text);
    }
  })
}

function sendRequestTest(url, data) {
  // send request
  $.ajax({
    url: url, // url where to submit the request
    type: "POST", // type of action POST || GET
    // dataType: "json", // data type
    contentType: "application/x-www-form-urlencoded",
    data: data, // post data || get data
    success: function(result) {
      // you can see the result from the console
      // tab of the developer tools
      alert(result);
      console.log(result);
      $(".result").html(result);
    },
    error: function(xhr, resp, text) {
      console.log(xhr, resp, text);
    }
  })
}
