function get_next() {
    $.getJSON("/review/next",
        function (data, textStatus, jqXHR) {
            if (textStatus == "success") {
                $.each(data, function (indexInArray, valueOfElement) { 
                    var img = $('<img></img>');
                    img.attr('src',valueOfElement.url);
                    img.addClass('img-responsive center-block');
                    var pic = $('<div></div>');
                    pic.addClass('item col-lg-2 col-md-3 col-sm-6'); 
                    pic.data('id',valueOfElement.id);
                    pic.append(img);
                    var pics = $('#pics');
                    pics.append(pic);
                    $('#pics').masonry('addItems',pic);
                });
                $('#pics').imagesLoaded(function(){
                    $('#pics').masonry('layout');
                });
            }
        }
    );
}
function get_next_pics(page) {
    $.getJSON("/page/"+page,
        function (data) {
            $.each(data, function (indexInArray, valueOfElement) { 
                var img = $('<img></img>');
                img.attr('src',valueOfElement.url);
                img.addClass('img-responsive center-block');
                var a = $('<a></a>');
                a.attr('href','/detail/'+valueOfElement.id);
                a.attr('target','_blank');
                a.append(img);
                var pic = $('<div></div>');
                pic.addClass('item col-lg-2 col-md-3 col-sm-6'); 
                //pic.css('margin-top','5px');
                pic.data('id',valueOfElement.id);
                pic.append(a);
                var pics = $('#pics');
                pics.append(pic);
                $('#pics').masonry('addItems',pic);
            });
            $('#pics').imagesLoaded(function(){
                $('#pics').masonry('layout');
            });
        }
    );

    current+=1;
}