function get_next() {
    $.getJSON("/review/next",
        function(data, textStatus, jqXHR) {
            if (textStatus == "success") {
                $('#pics').empty();
                $.each(data, function(indexInArray, valueOfElement) {
                    var img = $('<img>');
                    img.attr('src', valueOfElement.url);
                    img.addClass('img-responsive center-block');
                    var pic = $('<div>');
                    pic.addClass('item col-lg-2 col-md-3 col-sm-6');
                    pic.data('id', valueOfElement.id);
                    pic.append(img);
                    var pics = $('#pics');
                    pics.append(pic);
                    $('#pics').masonry('addItems', pic);
                });
                $('#pics').imagesLoaded(function() {
                    $('#pics').masonry('layout');
                });
            }
        }
    );
}

function get_next_pics(page) {
    current += 1;
    if (page % 5 == 0) {
        pics.empty();
    }
    $.getJSON("/page/" + page,
        function(data) {
            $.each(data, function(indexInArray, valueOfElement) {
                var img = $('<img>');
                img.attr('src', valueOfElement.url);
                img.addClass('img-responsive center-block');
                var pic = $('<div>');
                pic.addClass('item col-lg-2 col-md-3 col-sm-6');
                pic.data('id', valueOfElement.id);
                pic.append(img);
                img.imagesLoaded(function() {
                    $('#pics').masonry('layout');
                });
                $('#pics').append(pic);
                $('#pics').masonry('addItems', pic);
            });
        }
    );
}