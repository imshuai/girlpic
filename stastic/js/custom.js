$(function() {
    //Event for trial
    $('#NewTrialBtn').click(function() {
        $.getJSON("/try/new", function(json) {
            var tbody = $('#trial-list tbody:first');
            var tr = $('<tr></tr>');
            if (json.is_used == "是") {
                tr.addClass("danger");
            } else {
                tr.addClass("success");
            }
            tr.html("<td>" + json.token + "</td>\n<td>" + json.node + "</td>\n<td>" + json.expire + "</td>\n<td>" + json.is_used + "</td>");
            tbody.append(tr);
        });
    });


    //Event for bill
    $("#modalconfirm").on("shown.bs.modal", function(event) {
        var pId = $(event.relatedTarget).data("pid");
        $(this).find("#submitnewbill").attr("data-pid", pId);
    });

    $('#modalconfirm').on('hidden.bs.modal', function() {
        $('#nickname').val('');
    });

    $('#submitnewbill').click(function() {
        nickname = $('#nickname').val();
        if (nickname == '') {
            nickname = '匿名用户'
        }
        $('#modalconfirm').modal('hide');
        var pId = $(this).data('pid');
        // $('#modalcheckpaystatus').modal({
        //     backdrop: 'static',
        //     keyboard: false,
        //     remote: "/bill/new/"+pId
        // });
        $('#modalcheckpaystatus').find('div.modal-body').load("/bill/new/" + pId, function(responseTxt, statusTxt, xhr) {
            if (statusTxt == 'success') {
                $('#modalcheckpaystatus').modal({
                    backdrop: 'static',
                    keyboard: false
                });
            }
        });

    });

    // $('#modalcheckpaystatus').on('show.bs.modal', function(){
    //     var pId = $('#submitnewbill').data("pid");
    //     $(this).find('div.modal-content').load("/bill/new/"+pId);
    // });

    // $('#modalcheckpaystatus').on('loaded.bs.modal', function(){
    //     var bNo = $('#modalcheckpaystatus div.info').data('bno');
    //     console.log(bNo);
    //     sw = false;
    //     checkPayStatus(bNo);
    // });
    $('#modalcheckpaystatus').on('shown.bs.modal', function() {
        bNo = $('#modalcheckpaystatus div.info').data('bno');
        sw = false;
        checkPayStatus();
    });

    $('#modalcheckpaystatus').on('hidden.bs.modal', function() {
        sw = true;
        $('#nickname').val('');
    });

});

var sw;
var bNo;
var nickname;

function checkPayStatus() {
    if (sw) {
        return
    } else {
        // $.ajax({
        //     url: "/bill/query/"+bno,
        //     type: "get",
        //     async: false,
        //     cache: false,

        // });
        var span = $('#modalcheckpaystatus div.info li.li-info:last span');
        span.load('/bill/query/' + bNo, function(responseTxt, statusTxt, xhr) {
            if (span.text() == "支付状态：交易完成") {
                sw = true;
                $('#modalcheckpaystatus').modal('hide');
                $('#modalnewuser').find('div.modal-body').load("/user/new/" + bNo + '?name=' + nickname, function(responseTxt, statusTxt, xhr) {
                    if (statusTxt == 'success') {
                        $('#modalnewuser').modal({
                            backdrop: 'static',
                            keyboard: false
                        });
                    }
                });
            } else {
                setTimeout('checkPayStatus()', 3000);
            }
        });

    }
}