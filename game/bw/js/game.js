var ua = window.navigator.userAgent;
var de = document.documentElement;
var bwweb_config = {
	'default': {
		'activexobject_name': 'ZhushenPlugin.IEControlEx.5',
		'npplugin_name': 'npzhushenplugin',
		'npplugin_description': 'npzhushenplugin 1, 0, 5, [0-9]+',
		'classid': 'CAF0E663-05DA-40ea-A7EC-DA66292AD0C6',
		'nptype': 'application/ZhushenGame-plugin'
	}
};
var config = bwweb_config[zoneID] || bwweb_config['default'];
var activexobject_name = config.activexobject_name;
var npplugin_name = config.npplugin_name;
var npplugin_description = config.npplugin_description;
var iecontrol_object = '<object id="IEControl" classid="CLSID:' + config.classid + '" style="width:100%;height:100%">' + param + "</object>";
var npcontrol_object = '<object id="NPControl" type="' + config.nptype + '" style="width:100%;height:100%">' + param + "</object>";
var pt = 3, pr = 3, pb = 3, pl = 3;
function updateControlMarginTop(e){
  e.style.marginTop = "0";
}
function scale(e, pt, pr, pb, pl){
  if (e) {
    e.style.width = (de.clientWidth - pl - pr) + "px";
    e.style.height = (de.clientHeight - pt - pb) + "px";
    updateControlMarginTop(e);
  }
}
function setResolution(e, w, h){
  if (e) {
    w = w || de.clientWidth;
    h = h || de.clientHeight;
    e.style.width = w + "px";
    e.style.height = h + "px";
    updateControlMarginTop(e);
  }
}
function fullScreen(){
  var e = document.getElementById("IEControl") || document.getElementById("NPControl");
  scale(e, pt, pr, pb, pl);
}
function resize(e){
  fullScreen();
}
function checkPlugin(){
  if (window.ActiveXObject) {
    var obj = null;
    try {
      obj = new ActiveXObject(activexobject_name);
    } catch(e) {
      obj = null;
    }
    return (obj != null);
  } else {
    navigator.plugins.refresh(false);
    return (navigator.plugins[npplugin_name] !== undefined && (new RegExp(npplugin_description)).test(navigator.plugins[npplugin_name].description));
  }
}
function loopScale(e, id, n, func, t){
  if (typeof(n) === "number" && n > 0) {
    e = e || document.getElementById(id);
    if (e) {
      func(e);
    }
    setTimeout(function(){ loopScale(e, id, n - 1, func, t * 2); }, t);
  }
}
function startGameObject(){
  replacePage();
  disableKey();
  windowProcess();
  var control, id;
  if (window.ActiveXObject) {
    document.getElementById("plugin").innerHTML = iecontrol_object;
    id = "IEControl";
    control = document.getElementById("IEControl");
    try{ IEControl.SetJsCallback(control_callback); }catch(e){}
  } else {
    document.getElementById("plugin").innerHTML = npcontrol_object;
    id = "NPControl";
    control = document.getElementById("NPControl");
    try{ NPControl.JsCallback = control_callback; }catch(e){}
  }
  if (control) {
    resize(control);
    control.onload = function(){ resize(control); };
    window.onload = function(){ resize(control); };
    window.onresize = function(){ resize(control); };
  } else {
    loopScale(control, id, 3, function(e){ resize(control); }, 200);
  }
}
function loopCheck(){
  if (checkPlugin())
    startGameObject();
  else
    setTimeout(loopCheck, 500);
}
function setUp(){
  document.getElementById("setuphref").href = setuppath;
  document.getElementById("start").style.display = "block";
}
function replacePage(){
  document.getElementById("start").style.display = "none";
}
function readyForStartGame(){
  if (checkPlugin()) {
    startGameObject();
  } else {
    setUp();
    setTimeout(loopCheck, 3000);
  }
}
jQuery(document).ready(function(){
  readyForStartGame();
});

function disableKey(){
  $("body").bind("keydown keypress", function(){
    // Alt + Up,Down,Left,Right
    if((window.event.altKey) && ((window.event.keyCode==37) || (window.event.keyCode==39))) {
      event.returnValue=false;
    }
    // Backspace, F5
    else if ((event.keyCode==8) || (event.keyCode==116)) {
      event.keyCode=0;
      event.returnValue=false;
    }
    // Ctrl + N
    else if ((event.ctrlKey) && (event.keyCode==78)) {
      event.returnValue=false;
    }
    // Shift + F10
    else if ((event.shiftKey) && (event.keyCode==121)) {
      event.returnValue=false;
    }
  });
}
function windowProcess(){
  window.onbeforeunload = function(e) {
    var e = e || window.event;
    var msg = gameName + '正在运行，确定要离开吗？';
    if (e) e.returnValue = msg;
    return msg;
  }
}
function back(){
  window.location.href = loginurl;
}
function safeReturnSrvSelect(){
  window.onbeforeunload = null;
  back();
}
function control_callback(msg){
  if (msg == "ReturnToSrvSelect") {
    safeReturnSrvSelect();
  }
}
