function get_next() {
    $.getJSON("/review/next",
        function(data, textStatus, jqXHR) {
            if (textStatus == "success") {
                $('#pics').masonry('remove', $('#pics>div'));
                $(window).scrollTop(0);
                $.each(data, function(indexInArray, valueOfElement) {
                    var img = $('<img>');
                    img.attr('src', valueOfElement.url);
                    img.addClass('img-responsive center-block');
                    var pic = $('<div>');
                    pic.addClass('item col-lg-2 col-md-3 col-sm-6');
                    pic.data('id', valueOfElement.id);
                    img.imagesLoaded(function() {
                        $('#pics').masonry('layout');
                    });
                    pic.append(img);
                    $('#pics').append(pic);
                    $('#pics').masonry('addItems', pic);
                });
            }
        }
    );
}

function get_next_pics(page) {
    current += 1;
    if (page % 5 == 0) {
        $('#pics').masonry('remove', $('#pics>div'));
        $(window).scrollTop(0);
    }
    $.getJSON("/page/" + page,
        function(data) {
            $.each(data, function(indexInArray, valueOfElement) {
                var img = $('<img>');
                img.attr('src', valueOfElement.url);
                img.addClass('img-responsive center-block');
                var a = $('<a>');
                a.attr('href', valueOfElement.url).attr('target', "_blank").append(img);
                var pic = $('<div>');
                pic.addClass('item col-lg-2 col-md-3 col-sm-6');
                pic.data('id', valueOfElement.id);
                pic.append(a);
                img.imagesLoaded(function() {
                    $('#pics').masonry('layout');
                });
                $('#pics').append(pic);
                $('#pics').masonry('addItems', pic);
            });
        }
    );
}