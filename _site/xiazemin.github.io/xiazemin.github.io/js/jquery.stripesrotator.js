/*
 * Stripes Rotator - jQuery plugin for rotating images with stripes.
 *
 * Copyright (c) 2010 Giulio Iotti
 *
 * Licensed under the MIT license:
 *   http://www.opensource.org/licenses/mit-license.php
 *
 */
(function($) {
	$.fn.stripesRotator = function(options) {
		var root = this;		/* Must be called from a div element that will be filled with the images. */
		var currImgIndex = 1;

		var defaults = {
			width: 15,			/* Width of each stripe. */
			total: 0,			/* Total number of stripes. (Calculated automatically.) */
			images: null,		/* Array of images (as jQuery objects). */
			pause: 0,			/* Time to pause between cycles. */
			stripeTime: 1200,	/* Time to animate a single stripe. */
			waitTime: 200,		/* Time to wait before animating another stripe. */
			blockId: 'srParent' /* Id of the block to be temporarily created. */
		};

		var data = $.extend({}, defaults, options);
		var images = [];	/* Keeps a vector of jQuery objects representing each image. */

		function pushAndAnimate(i) {
			images.push($(this));
			
			if (images.length > 1) {
				/* If we have enough images, start the animation. */
				_animate();
			}
		}

		/* Get all images inside an element and returns them as an array of jQuery objects. */
		function getAll(parent) {
			parent.children('img').each(pushAndAnimate);
		}

		function changeBlock(i) {
			if (i >= data.total) {
				return;
			}

			$('#'+data.blockId+i).animate({ scrollLeft: data.width+'px' }, data.stripeTime);

			/* Quirk to make setTimeout call the function itself. */
			var selfCallback = function() {
				changeBlock(i+1);
			};

			window.setTimeout(selfCallback, data.waitTime);
		}

		function doAnimation(imgSrc) {
			var firstImgSrc = imgSrc.attr('src');
			var height = imgSrc.attr('height');

			if (data.total <= 0) {
				data.total = (parseInt(imgSrc.attr('width'), 10) / data.width) + 1;
			}

			var argsPar = {
				'height': height + 'px',
				'width': data.width + 'px',
				'overflow': 'hidden',
				'position': 'absolute',
				'margin-left': '0px'
			};

			var argsImgCont = {
				'position': 'absolute',
				'height': height + 'px',
				'width': (data.width * 3) + 'px',
				'overflow': 'hidden'
			};

			data.currImg = imgSrc;

			for (i = 0; i < data.total; i++) {
				var imgParent = $(document.createElement('div'));
				var imageContainer = $(document.createElement('div'));
				var img = $(document.createElement('img'));
				var leftMargin = 0;

				if (argsPar['margin-left']) {
					leftMargin = parseInt(argsPar['margin-left'].replace('px', ''), 10);
				}

				argsImgCont['margin-left'] = (data.width * i - leftMargin) + 'px';
				argsPar['margin-left'] = (data.width * i) + 'px';

				imgParent.attr('id', data.blockId + i);
				imgParent.css(argsPar);

				img.attr('src', firstImgSrc);

				imageContainer.css(argsImgCont);
				imageContainer.append(img);

				var pos = -(data.width * i);

				if (pos === 0) {
					img.css({ 'margin-left': data.width + 'px' });
				}

				img.css({
					'top': 0,
					'left': pos + 'px',
					'position': 'absolute'
				});
				imageContainer.scrollLeft(data.width * i);

				imgParent.append(imageContainer);
				root.append(imgParent);
			}

			/* Do the camaleon. */
			changeBlock(0);

			/* Cleanup after the party. */
			window.setTimeout(function(img1) {
				root.css('background-image', 'url(' + data.currImg.attr('src') + ')');

				root.children('div').each(function(i) {
					$(this).remove();
				});

				_animate();
			}, data.waitTime * data.total + data.stripeTime + data.pause);
		}
		
		function _animate() {
			if (currImgIndex >= images.length) {
				currImgIndex = 0;
			}

			doAnimation(images[currImgIndex]);
			currImgIndex++;
		}
		
		function loadImages(str_images) {
			for (var i = 0; i < str_images.length; i++) {
				if (typeof(str_images[i])=='string') { // It's not a jQuery object;
					
					$("<img />").attr('src', str_images[i]).bind("load", pushAndAnimate);
				}
			}
		}
		
		/* It is possible to specify a jQuery object (for example a hiddden div)
		   that contains various images. */
		if (data.images && data.images.length > 1) {
			loadImages(data.images);
		} else {
			getAll(data.images);
		}

		/* Since we don't do any useful modification to the object itself,
		   we just pass it on to the next action in the chain. */
		return this;
	};
})(jQuery);
