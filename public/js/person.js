$(function() {
    $('#my-tab-content').delegate('button[name=check]', 'click', function() {
        $.ajax({  
         url : "/getform",    
         data : {id:$(this).val()},
         type : "post",  
         cache : false,  
         dataType : "json",  
         success:getform   
         }); 
        $('#id').val($(this).val());
    });
    
    $('#my-tab-content').delegate('button[name=delete]', 'click', function() {
        $.ajax({  
         url : "/deleteform",    
         data : {id:$(this).val()},
         type : "post",  
         cache : false,  
         dataType : "json",  
         success: deleteform  
         }); 
    });
    $('#delete').click(function() {
        $.ajax({  
         url : "/deleteform",    
         data : {id:$("#id").val()},
         type : "post",  
         cache : false,  
         dataType : "json",  
         success: deleteform  
         }); 
    });      
    
    $('#checked').click(function() {
        $.ajax({  
         url : "/updateform",    
         data : {name:$("#name").val(), gender:$("#gender").val(), birth:$("#birth").val(), believe:$("#believe").val(), unit:$("#unit").val(), mobile:$("#mobile").val(), email:$("#email").val(), positions:$("#positions").val(), director:$("#director").val(), evaluation:$("#evaluation").val(), id:$("#id").val(), type:2},
         type : "post",  
         cache : false,  
         dataType : "json",  
         success:  updateform 
         }); 
    });         
});
function updateform(json) {
    if (json.Status !== undefined) {
        $('#updateModal').modal('toggle');
        $('#updateModal p').text(json.Status + ":" + json.Text);
        location.reload();
        return
    }    
}
function deleteform(json) {
    if (json.Status !== undefined) {
        $('#deleteModal').modal('toggle');
        $('#deleteModal p').text(json.Status + ":" + json.Text);
        location.reload();
        return
    }    
}
function getform(json) {
        $('#myModal').modal('toggle');        
        $('#name').val(json.Name);
        $('#gender').val(json.Gender);
        $("#birth").val(json.Birth);
        $('#believe').val(json.Believe);
        $('#unit').val(json.Unit);
        $('#mobile').val(json.Mobile);
        $('#email').val(json.Email);
        $('#positions').val(json.Positions);
        $('#director').val(json.Director);
        $('#evaluation').val(json.Evaluation);
        
}
