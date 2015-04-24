

$(function() {
    $("#login").click(function() {
        $.ajax({  
         url : "/rLogin",    
         data : {username:$("#username").val(), password:$("#password").val()},
         type : "post",  
         cache : false,  
         dataType : "json",  
         success: callback   
         }); 
    });  
});
    

function callback(json) {
    if (json.Status !== 200) {
        $('#myModal').modal('toggle');
        $('#myModal p').text(json.Status + ":" + json.Text);
        return
    }
     window.location.assign("/regform/admin/team");
//	window.location.assign("/");
}

