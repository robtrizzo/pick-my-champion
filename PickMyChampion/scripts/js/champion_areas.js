$(document).ready(function() {
    loadPortraits();
    makeDraggable();
    makeDroppable();
});

function makeDraggable() {
    $(".championPortrait").draggable({
        scroll: false,
        revert: 'invalid',
        stack: false,
        cursor: "pointer"
    });
};

function makeDroppable() {
    $(".championDropArea").droppable({
        //accept: ".championPortrait",
        accept: function(elem) {
            if($(this).hasClass("hasChampion")){
                return false;
            }
            return true;
        },
        drop: function(event, ui) {
            var $this = $(this);

            $(this).addClass("hasChampion");
            $(ui.draggable).parent().removeClass("hasChampion")
            $(ui.draggable).appendTo($(this));
            ui.draggable.position({
                my: "center",
                at: "center",
                of: $this
            });
            var championId = ui.draggable.attr("id");
            var droppableId = $(this).attr("id");

            $.ajax({
                url: '/scripts/go/championDropped',
                type: 'post',
                dataType: 'text',
                data: {
                    champion_name: championId,
                    droppable_name: droppableId
                },
                success: function(response) {
                    console.log(response)
                }
            });
        }
    });
};

function loadPortraits() {
    dir = '/static/img/champion-portraits/';
    $.ajax({
        url: '/scripts/go/listDir',
        type: 'post',
        dataType: 'text',
        data: {
            dir_path: dir
        },
        success: function(response) {
            imgs = response.split(",");
            for (i = 0; i < imgs.length; i++) {
                loadImage(dir, imgs[i], i);
            }
            makeDroppable();
            makeDraggable();
        }
    });
};

function loadImage(dir, img_fn, champion_num) {
    var div = document.createElement("div");
    div.className = "championDropArea hasChampion";
    div.id = "champion_default" + champion_num
    var img = document.createElement("img");
    img.src = dir + "/" + img_fn;
    img.id = img_fn.substr(0, img_fn.lastIndexOf('.'));;
    img.className = "championPortrait";
    var src = document.getElementById("championList");
    div.appendChild(img);
    src.appendChild(div);

}