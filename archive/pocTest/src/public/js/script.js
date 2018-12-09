$(document).ready(function() {
  // do jQuery

  $("#dynamicField").submit(function(event) {
    event.preventDefault();
    console.log($("#dynamicFieldName").val(), $("#dynamicFieldVal").val());
    $('<input>').attr({
      // type: 'hidden',
      name: $("#dynamicFieldName").val(),
      value: $("#dynamicFieldVal").val()
    }).appendTo('#formFields');
  });

  $("#insertForm").submit(function(event) {
    console.log(event);

    // prepair data for request
    data = processData.bind(this).call();

    sendRequest("insert", data);

    event.preventDefault();
  });


  $("#updateForm").submit(function(event) {
    console.log(event);

    // prepair data for request
    data = processData.bind(this).call();

    sendRequest("update", data);

    event.preventDefault();
  });

  $("#validateForm").submit(function(event) {
    console.log(event);

    // prepair data for request
    data = processData.bind(this).call();

    sendRequest("validate", data);

    event.preventDefault();
  });

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
