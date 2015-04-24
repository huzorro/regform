$(function() {
    $("#submit").click(function() {
        $.ajax({  
         url : "/regform/person/submit",    
         data : {name:$("#name").val(), gender:$("#gender").val(), birth:$("#birth").val(), believe:$("#believe").val(), unit:$("#unit").val(), mobile:$("#mobile").val(), email:$("#email").val(), positions:$("#positions").val(), director:$("#director").val(), evaluation:$("#evaluation").val()},
         type : "post",  
         cache : false,  
         dataType : "json",  
         success: personSubmit   
         }); 
    });
});


function personSubmit(json) {
    if (json.Status !== undefined) {
        $('#myModal').modal('toggle');
        $('#myModal p').text(json.Status + ":" + json.Text);
        return
    }
}