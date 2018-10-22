$(document).ready(function(){
  let toggled = false
  $("#menu-toggle").click(function(e){
    let display = toggled ? "none" : "inline"
    let width = toggled ? "50px" : "200px"
    $("#left-bar").animate({width: width}, 400)
    $("#left-bar span").css("display", display)
    
    toggled = !toggled  
  })
})