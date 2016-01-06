$(function() {
  var MAX_DEPTH = 3;

  $("ol.categories").sortable({
    handle: 'i.icon-move',
    onDrop: function($item, container, _super, event) {
      determineSubcategoryAdditionAbility($item);
      _super($item, container);
    }
  });

  $(document).on('click', '#add-root-category', function() {
    addCategory();
  });

  $(document).on('click', '.remove-button', function() {
    var el = $(this).parent().parent();
    el.remove();
  });

  $(document).on('click', '.subcat-button', function() {
    var el = $(this).parent().parent().parent();
    var ol = el.children()[1];
    if(!ol) {
      ol = $("<ol />");
      el.append(ol);
    } else {
      ol = $(ol)
    }
    addCategory(ol, "Untitled Category", "#333");
  });

  $(document).on('blur', '.category-name-edit', function() {
    var $this = $(this);
    var content = $this.val();
    var spanEl = $('<span class=\"title\" />');
    spanEl.html(content);
    $this.replaceWith(spanEl);
  });

  $(document).on('click', '.title', function() {
    var $this = $(this);
    var content = $this.html();
    var inputEl = $('<input class="category-name-edit" name="category-name-edit" type="text">');
    inputEl.val(content);
    $this.replaceWith(inputEl);
    inputEl.onfocus = function() {
      inputEl.select();

      inputEl.onmouseup = function() {
        inputEl.onmouseup = null;
        return false;
      };
    };
    setTimeout(function() {
      inputEl.select();
    }, 1);
  });

  function createCategoryElement(title, color, shouldNest) {
    var newItemLi = $("<li />");
    var newItemDiv = $("<div class=\"category-item\"></div>");
    var titleSpan = $("<span class=\"title\" />");
    var dragHandle = $("<i class=\"icon-move glyphicon glyphicon-th-list\"></i>");
    var buttonsContainer = $("<div class=\"pull-right\"></div>");
    // if(shouldNest) {
    //   var addSubCategoryButton = $("<button type=\"button\" class=\"subcat-button btn btn-default btn-xs\"><span class=\"glyphicon glyphicon-plus\"></span><span class=\"hidden-xs\"> Sub Category</span></button>");
    // }
    var editButton = $("<button type=\"button\" class=\"edit-button btn btn-default btn-xs\"><span class=\"glyphicon glyphicon-edit\"></span><span class=\"hidden-xs\"> Edit</span></button>");
    var removeButton = $("<button type=\"button\" class=\"remove-button btn btn-default btn-xs\"><span class=\"glyphicon glyphicon-remove\"></span><span class=\"hidden-xs\"> Remove</span></button>");

    //buttonsContainer.append(addSubCategoryButton);
    buttonsContainer.append(editButton);
    buttonsContainer.append(removeButton);

    newItemDiv.append(dragHandle);
    titleSpan.html(title);
    newItemDiv.append(titleSpan);
    newItemDiv.append(buttonsContainer);
    newItemLi.append(newItemDiv);

    return newItemLi;
  }

  function addCategory(parentEl, title, color) {
    var maxDepth = 3
    var parent = parentEl || $("ol.categories");
    var elDepth = $(parent).parents('ol').length;
    var color = color || parent.css('background-color');
    var title = title || "Untitled Category";
    var newCat = createCategoryElement(title, color, canNestCategory(parent, MAX_DEPTH - 1));
    parent.prepend(newCat);
    determineSubcategoryAdditionAbility(newCat, MAX_DEPTH);
  }

  function determineSubcategoryAdditionAbility(el, depth) {
    var addSubCategoryButton = $("<button type=\"button\" class=\"subcat-button btn btn-default btn-xs\"><span class=\"glyphicon glyphicon-plus\"></span><span class=\"hidden-xs\"> Sub Category</span></button>");
    var hasAddSubcategoryButton = (el.find('.subcat-button').length > 0);

    if(canNestCategory(el, depth)) {
      if(hasAddSubcategoryButton) { return; }
      el.find(".pull-right").append(addSubCategoryButton);
    } else {
      if(!hasAddSubcategoryButton) { return; }
      el.find(".subcat-button").remove();
    }
  }

  function canNestCategory(categoryEl, maxDepth) {
    maxDepth = maxDepth || MAX_DEPTH;
    var elDepth = $(categoryEl).parents('ol').length;
    return elDepth < maxDepth;
  }
});