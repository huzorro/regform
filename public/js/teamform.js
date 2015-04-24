$(function() {
    $("#submit").click(function() {
        $.ajax({  
         url : "/regform/team/submit",    
         data : {university:$("#university").val(), level:$("#level").val(), leader:$("#leader").val(), leader_gender:$("#leader_gender").val(), leader_contact:$("#leader_contact").val(), pm:$("#pm").val(), pm_gender:$("#pm_gender").val(), pm_contact:$("#pm_contact").val(), topic:$("#topic").val(), intro:$("#intro").val()},
         type : "post",  
         cache : false,  
         dataType : "json",  
         success:teamSubmit
         }); 
    });
});


function teamSubmit(json) {
    if (json.Status !== undefined) {
        $('#myModal').modal('toggle');
        $('#myModal p').text(json.Status + ":" + json.Text);
        return
    }
}