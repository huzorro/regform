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
         data : {university:$("#university").val(), level:$("#level").val(), leader:$("#leader").val(), leader_gender:$("#leader_gender").val(), leader_contact:$("#leader_contact").val(), pm:$("#pm").val(), pm_gender:$("#pm_gender").val(), pm_contact:$("#pm_contact").val(), topic:$("#topic").val(), intro:$("#intro").val(), id:$("#id").val(), type:1},
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
        $('#university').val(json.University);
        $('#level').val(json.Level);
        $("#leader").val(json.Leader);
        $('#leader_gender').val(json.LeaderGender);
        $('#leader_contact').val(json.LeaderContact);
        $('#pm').val(json.Pm);
        $('#pm_gender').val(json.PmGender);
        $('#pm_contact').val(json.PmContact);
        $('#topic').val(json.Topic);
        $('#intro').val(json.Intro);
        
}
