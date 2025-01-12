package main

import (
	"fmt"
	"github.com/fouched/social/internal/repo"
	"log"
	"math/rand"
)

//TODO move this to a lib that can be included or use https://github.com/jaswdr/faker
// this file slows down initial compile time

var userNames = []string{
	"aaliyah", "aaron", "abigail", "adam", "adrian", "aiden", "ainsley", "aisha", "alan", "alexa",
	"alexander", "alexis", "alice", "alicia", "alison", "alyssa", "amber", "amelia", "amina", "amy",
	"andrew", "angela", "anita", "anna", "anthony", "april", "ariana", "arianna", "ashley", "ashton",
	"aubrey", "audrey", "austin", "ava", "avery", "barbara", "benjamin", "beth", "beverly", "blake",
	"bobby", "brandon", "breanna", "brenda", "brian", "briana", "brianna", "brittany", "brooke", "bryan",
	"caleb", "cameron", "cara", "carlos", "carly", "carmen", "caroline", "carter", "casey", "cassandra",
	"catherine", "cecilia", "chad", "charles", "charlotte", "chase", "cheryl", "chloe", "chris", "christina",
	"christine", "christopher", "cindy", "claire", "clara", "clarence", "claudia", "clayton", "cody", "colin",
	"connor", "courtney", "craig", "crystal", "curtis", "cynthia", "daisy", "dakota", "dana", "daniel",
	"danielle", "david", "dean", "deanna", "debbie", "deborah", "denise", "dennis", "derek", "derrick",
	"desiree", "diana", "diane", "dominic", "donald", "donna", "doris", "dorothy", "dylan", "earl",
	"edgar", "eddie", "eden", "edith", "edward", "elaine", "elena", "elijah", "elizabeth", "ella",
	"ellen", "elliot", "emily", "emma", "eric", "erica", "erik", "erin", "ethan", "eugene",
	"eva", "evan", "evelyn", "faith", "felicia", "felix", "fernando", "fiona", "frances", "francisco",
	"frank", "gabriel", "gail", "gareth", "gary", "gavin", "george", "georgia", "gerald", "geraldine",
	"gina", "glen", "glenn", "gloria", "grace", "greg", "gregory", "gwen", "hailey", "hannah",
	"harold", "harry", "hazel", "heather", "heidi", "helen", "henry", "holly", "howard", "ian",
	"irene", "isaac", "isabel", "isabella", "jack", "jackie", "jackson", "jacob", "jacqueline", "jade",
	"jake", "james", "jamie", "jane", "janet", "janice", "jared", "jason", "jasmine", "jay",
	"jayden", "jean", "jeff", "jeffery", "jeffrey", "jen", "jenna", "jennifer", "jenny", "jeremy",
	"jerome", "jerry", "jesse", "jessica", "jill", "joan", "joanna", "joe", "joel", "john",
	"johnny", "jonathan", "jordan", "jorge", "jose", "joseph", "josh", "joshua", "joy", "joyce",
	"juan", "judith", "judy", "julia", "julian", "julie", "june", "justin", "kaitlyn", "kara",
	"karen", "kari", "karla", "kate", "katelyn", "katherine", "kathleen", "kathryn", "kathy", "katie",
	"katrina", "kay", "kayla", "keith", "kelly", "ken", "kendra", "kenneth", "kevin", "kim",
	"kimberly", "kirk", "krista", "kristen", "kristin", "kristina", "kyle", "kylie", "lance", "larry",
	"laura", "lauren", "leah", "lee", "leon", "leonard", "leslie", "liam", "lillian", "lindsay",
	"lisa", "logan", "lori", "lorraine", "louis", "lucas", "luis", "lydia", "lynn", "mackenzie",
	"madeline", "madison", "maggie", "mallory", "mandy", "marc", "marcia", "margaret", "maria", "mariah",
	"marilyn", "mario", "marissa", "mark", "marlene", "marsha", "martha", "martin", "marvin", "mary",
	"mason", "matt", "matthew", "maureen", "megan", "melanie", "melissa", "melvin", "meredith", "michael",
	"michaela", "michelle", "miguel", "mike", "mindy", "miranda", "misty", "mitch", "mitchell", "molly",
	"monica", "morgan", "nancy", "naomi", "natalie", "natasha", "nathan", "nathaniel", "neal", "ned",
	"neil", "nelson", "nicholas", "nicole", "nina", "noah", "noelle", "norma", "oliver", "olivia",
	"omar", "oscar", "paige", "pam", "pamela", "patricia", "patrick", "paul", "paula", "pauline",
	"peggy", "penny", "peter", "phil", "philip", "phyllis", "rachel", "ralph", "randy", "ray",
	"raymond", "rebecca", "regina", "renee", "rex", "rhonda", "richard", "rick", "rita", "rob",
	"robert", "roberta", "robin", "rodney", "roger", "ron", "ronald", "rosa", "rose", "rosemary",
	"ross", "roy", "ruben", "ruby", "russell", "ruth", "ryan", "sabrina", "sally", "samantha",
	"samuel", "sandra", "sara", "sarah", "scott", "sean", "shane", "shannon", "sharon", "shaun",
	"sheila", "shelby", "sherry", "shirley", "sidney", "sierra", "simon", "sofia", "sonia", "sophia",
	"spencer", "stacey", "stacy", "stanley", "stefanie", "stephanie", "stephen", "steve", "steven", "sue",
	"summer", "susan", "susie", "suzanne", "sydney", "sylvia", "tamara", "tami", "tammy", "tanya",
	"tara", "taylor", "ted", "teresa", "terri", "terry", "theodore", "theresa", "thomas", "tiffany",
	"tina", "todd", "tom", "tommy", "toni", "tony", "tracey", "tracy", "travis", "trent",
	"trevor", "tricia", "tristan", "troy", "tyler", "valerie", "vanessa", "vera", "veronica", "vicki",
	"vicky", "victor", "victoria", "vincent", "violet", "virginia", "vivian", "wade", "walter", "wanda",
	"warren", "wayne", "wendy", "wesley", "whitney", "william", "willie", "wyatt", "xander", "xavier",
	"yasmin", "yolanda", "zachary", "zoe",
}

var domainNames = []string{
	"microsoft.com", "gmail.com", "yahoo.co.uk", "tech.net", "amazon.com", "apple.com", "facebook.com", "twitter.com", "linkedin.com", "instagram.com",
	"netflix.com", "github.com", "reddit.com", "quora.com", "medium.com", "slack.com", "zoom.us", "dropbox.com", "salesforce.com", "adobe.com",
	"oracle.com", "ibm.com", "intel.com", "nvidia.com", "spotify.com", "pinterest.com", "tumblr.com", "flickr.com", "snapchat.com", "tiktok.com",
	"whatsapp.com", "telegram.org", "messenger.com", "wechat.com", "paypal.com", "stripe.com", "shopify.com", "etsy.com", "ebay.com", "walmart.com",
	"target.com", "bestbuy.com", "homedepot.com", "lowes.com", "macys.com", "costco.com", "ikea.com", "wayfair.com", "kohls.com", "samsclub.com",
	"gap.com", "oldnavy.com", "bananarepublic.com", "jcrew.com", "nordstrom.com", "bloomingdales.com", "saksfifthavenue.com", "neimanmarcus.com", "anthropologie.com", "urbanoutfitters.com",
	"freepeople.com", "zara.com", "hm.com", "asos.com", "boohoo.com", "forever21.com", "shein.com", "romwe.com", "zaful.com", "fashionnova.com",
	"aliexpress.com", "jd.com", "rakuten.com", "taobao.com", "tencent.com", "baidu.com", "weibo.com", "qq.com", "sohu.com", "163.com",
	"xinhuanet.com", "chinadaily.com.cn", "people.com.cn", "sina.com.cn", "cntv.cn", "ifeng.com", "sogou.com", "alipay.com", "didi.com", "meituan.com",
	"xiaomi.com", "huawei.com", "lenovo.com", "oppo.com", "vivo.com", "bytedance.com", "kuaishou.com", "douyu.com", "huya.com", "bilibili.com",
	"jd.com", "coupang.com", "naver.com", "daum.net", "kakao.com", "line.me", "mercari.com", "rakuma.com", "shopee.com", "lazada.com",
	"tokopedia.com", "bukalapak.com", "gojek.com", "grab.com", "sephora.com", "ulta.com", "glossier.com", "nyxcosmetics.com", "jcpenney.com", "bonton.com",
	"lordandtaylor.com", "barneys.com", "bergdorfgoodman.com", "saks.com", "nordstromrack.com", "bhphotovideo.com", "adorama.com", "jet.com", "overstock.com", "cb2.com",
	"westelm.com", "potterybarn.com", "crateandbarrel.com", "worldmarket.com", "pier1.com", "havertys.com", "ashleyfurniture.com", "roomstogo.com", "raymourflanigan.com", "mgbwhome.com",
	"article.com", "joybird.com", "insideweather.com", "burrow.com", "thuma.co", "floydhome.com", "avocado.com", "casper.com", "purple.com", "helixsleep.com",
	"nectar.com", "tuftandneedle.com", "leesa.com", "saatva.com", "amerisleep.com", "sealy.com", "tempurpedic.com", "simmons.com", "beautyrest.com", "steinhafels.com",
	"valuecityfurniture.com", "bassettfurniture.com", "stickley.com", "thomasville.com", "ethanallen.com", "boh.com", "horchow.com", "perigold.com", "allmodern.com", "scandinaviandesigns.com",
	"lapensedefurniture.com", "maisoncorbeil.com", "mobilia.com", "eq3.com", "arhaus.com", "roomandboard.com", "designtoscano.com", "kensingtonrowcollection.com", "grandinroad.com", "ballarddesigns.com",
	"brooksbrothers.com", "ralphlauren.com", "vineyardvines.com", "madewell.com", "tommybahama.com", "lillypulitzer.com", "talbots.com", "anntaylor.com", "loft.com", "loftoutlet.com",
	"thelimited.com", "anntaylorloft.com", "hollisterco.com", "abercrombie.com", "abercrombiekids.com", "aeropostale.com", "americaneagle.com", "aerie.com", "athleta.com", "paige.com",
	"hudsonjeans.com", "joesjeans.com", "dl1961.com", "goodamerican.com", "nydj.com", "levi.com", "wrangler.com", "lee.com", "truefit.com", "hm.com",
	"zara.com", "pullandbear.com", "bershka.com", "stradivarius.com", "mango.com", "desigual.com", "esprit.com", "superdry.com", "jackwills.com", "bench.com",
	"jadedldn.com", "hugoboss.com", "armani.com", "versace.com", "gucci.com", "dolcegabbana.com", "prada.com", "fendi.com", "bottegaveneta.com", "balenciaga.com",
	"saintlaurent.com", "chloe.com", "givenchy.com", "celine.com", "loewe.com", "michaelkors.com", "coach.com", "furla.com", "longchamp.com", "goyard.com",
	"lv.com", "dior.com", "chanel.com", "hermes.com", "burberry.com", "ralphlauren.com", "tommy.com", "hillfiger.com", "nautica.com", "calvinklein.com",
	"dkny.com", "diesel.com", "gstar.com", "guess.com", "bebe.com", "bcbg.com", "karenmillen.com", "tedbaker.com", "reiss.com", "warehouse.co.uk",
	"coast-stores.com", "phase-eight.com", "whistles.com", "allsaints.com", "muji.com", "universalstandard.com", "everlane.com", "thoughtclothing.com", "sezane.com", "aritzia.com",
	"reformation.com", "babaa.es", "laredoute.com", "seasaltcornwall.com", "bodenusa.com", "jigsaw-online.com", "poetryfashion.com", "thefoldlondon.com", "hobbs.com", "finerylondon.com",
	"aligne.co", "meandem.com", "ba-sh.com", "sandqvist.com", "swedishhasbeens.com", "ganni.com", "arket.com", "cosstores.com", "weekday.com", "monki.com",
	"stories.com", "clarks.co.uk", "timberland.com", "ugg.com", "drmartens.com", "newbalance.com", "converse.com", "vans.com", "skechers.com", "stevemadden.com",
	"aldoshoes.com", "shoecarnival.com", "dsw.com", "zappos.com", "onlineshoes.com", "shoemall.com", "shoebuy.com", "famousfootwear.com", "footlocker.com", "finishline.com",
	"champssports.com", "eastbay.com", "footaction.com", "ycmc.com", "tactics.com", "surfdome.com", "revolution.co.uk", "elementbrand.com", "quiksilver.com", "ripcurl.com",
	"billabong.com", "volcom.com", "oneill.com", "roxy.com", "hurley.com", "reebok.com", "adidas.com", "nike.com", "puma.com", "underarmour.com",
	"columbia.com", "patagonia.com", "northface.com", "arcteryx.com", "marmot.com", "jackwolfskin.com", "fjallraven.com", "rab.equipment", "haglofs.com", "berghaus.com",
	"montbell.com", "outdoorresearch.com", "trespass.com", "regatta.com", "hellyhansen.com", "mustangsurvival.com", "stearnsflotation.com", "gillmarine.com", "henrilloyd.com", "gul.com",
	"hansensurf.com", "usoutdoor.com", "evo.com", "bobsportschalet.com", "buyee.jp", "moosejaw.com", "backcountry.com", "rei.com", "sierratradingpost.com", "campmor.com",
	"sondors.com", "radpower.com",
}

var titles = []string{
	"The Art of Minimalism", "Travel Hacks for Budget Trips", "Healthy Eating on a Busy Schedule", "Top 10 Travel Destinations", "Mastering the Work-Life Balance",
	"Secrets to Successful Blogging", "DIY Home Decor Ideas", "Tips for Remote Working", "The Benefits of Meditation", "Understanding Cryptocurrency",
	"A Guide to Self-Care", "Healthy Recipes for Beginners", "The Future of Technology", "Productivity Tips for Entrepreneurs", "Exploring Local Cuisine",
	"Improving Your Mental Health", "Sustainable Living Practices", "How to Start a Podcast", "The Magic of Mindfulness", "Creative Writing Tips",
	"Fitness Routines for Busy People", "How to Grow a Garden", "Photography Tips for Beginners", "Effective Time Management Strategies", "Best Books for Entrepreneurs",
	"How to Save Money", "The Benefits of Yoga", "Creative Content Ideas", "Eco-Friendly Lifestyle Choices", "Social Media Marketing Tips",
	"How to Build a Personal Brand", "Travel Guide to Paris", "Cooking for One: Tips and Tricks", "Understanding Mental Health", "Financial Planning for Millennials",
	"The Importance of Sleep", "Home Workout Routines", "How to Stay Motivated", "The Power of Positive Thinking", "Budget-Friendly Meal Ideas",
	"The Science of Happiness", "How to Create a Vision Board", "Healthy Snack Ideas", "The Benefits of Journaling", "How to Declutter Your Space",
	"Setting Realistic Goals", "The Power of Gratitude", "Tech Gadgets You Need", "How to Improve Your Focus", "The Benefits of Reading",
	"Traveling Solo: Tips and Advice", "Healthy Habits to Adopt", "How to Reduce Stress", "The Importance of Hydration", "How to Create a Budget",
	"Best Apps for Productivity", "Mindful Eating Practices", "Tips for a Successful Career", "The Importance of Exercise", "How to Stay Organized",
	"The Benefits of Volunteering", "How to Manage Your Time", "Healthy Morning Routines", "How to Start a Side Hustle", "Traveling with Kids: Tips and Tricks",
	"Creative Instagram Story Ideas", "How to Boost Your Confidence", "The Benefits of Walking", "How to Create a Capsule Wardrobe", "The Power of Networking",
	"How to Practice Gratitude", "The Importance of Self-Discipline", "Creative Ways to Stay Active", "How to Plan a Staycation", "The Benefits of Minimalist Living",
	"How to Improve Your Sleep", "Healthy Lunch Ideas", "The Power of Positive Affirmations", "How to Create a Cozy Home", "Tips for Starting a Blog",
	"The Benefits of Journaling", "How to Create a Balanced Diet", "The Importance of Self-Care", "How to Stay Focused", "The Benefits of Meditation",
	"How to Reduce Waste", "The Power of Positive Habits", "How to Stay Fit", "Healthy Dinner Ideas", "How to Create a Routine",
	"The Importance of Hydration", "How to Manage Stress", "The Benefits of Exercise", "How to Stay Productive", "Healthy Snack Ideas",
	"The Power of Positive Thinking", "How to Create a Vision Board", "The Benefits of Reading", "How to Stay Organized", "Healthy Habits to Adopt",
}

var contents = []string{
	"Living a minimalist lifestyle involves reducing your possessions to only the essentials. By decluttering your home and life, you create a more peaceful and focused environment. Minimalism encourages you to evaluate what truly brings you joy and eliminate what doesn't. This can lead to increased productivity, reduced stress, and a greater appreciation for the things you do have. It's not just about having fewer things, but also about making intentional choices that align with your values and goals. Embracing minimalism can help you live a more meaningful and fulfilling life.",
	"Traveling on a budget doesn't mean you have to miss out on amazing experiences. With a little planning and some smart strategies, you can explore the world without breaking the bank. Start by looking for cheap flights and accommodations. Websites like Skyscanner and Airbnb can help you find the best deals. Consider traveling during off-peak times to save money on flights and hotels. Take advantage of free activities and attractions, such as hiking, walking tours, and visiting museums on free admission days. Eating like a local and using public transportation can also help you save money while experiencing the culture of your destination.",
	"Maintaining a healthy diet can be challenging when you have a busy schedule, but it's not impossible. Start by planning your meals ahead of time and prepping ingredients in advance. This can save you time and make it easier to make healthy choices. Focus on simple, nutrient-dense foods like fruits, vegetables, whole grains, and lean proteins. Keep healthy snacks on hand, such as nuts, yogurt, and fresh fruit, to avoid reaching for unhealthy options when you're hungry. Remember to stay hydrated and make time for regular meals to keep your energy levels up throughout the day. With a little effort, you can stay healthy even when you're short on time.",
	"The work-life balance is something many people struggle with, but it's essential for maintaining your overall well-being. To achieve a better balance, start by setting clear boundaries between work and personal time. This might mean turning off work emails after a certain hour or designating specific times for family activities. Make sure to prioritize self-care and schedule regular breaks throughout the day to recharge. It's also important to communicate your needs with your employer and loved ones. By finding ways to manage your time effectively and making self-care a priority, you can achieve a healthier work-life balance and enjoy greater overall satisfaction.",
	"Meditation has been practiced for thousands of years and offers numerous benefits for both mental and physical health. By taking just a few minutes each day to sit quietly and focus on your breath, you can reduce stress, improve concentration, and increase feelings of well-being. Meditation can also help you develop greater self-awareness and emotional resilience. There are many different types of meditation, so it's important to find a practice that works for you. Whether you prefer guided meditations, mindfulness, or loving-kindness meditation, the key is to make it a regular part of your routine.",
	"Cryptocurrency is a digital or virtual form of currency that uses cryptography for security. Unlike traditional currencies, cryptocurrencies operate on a decentralized system called blockchain, which records all transactions across a network of computers. This makes them less susceptible to fraud and manipulation. Bitcoin, the first and most well-known cryptocurrency, was created in 2009. Since then, thousands of other cryptocurrencies have emerged, each with their unique features and uses. Cryptocurrencies offer several advantages, such as lower transaction fees and increased privacy. However, they also come with risks, including price volatility and regulatory uncertainty.",
	"Starting a podcast can be a great way to share your ideas and connect with a broader audience. To get started, you'll need to choose a topic that you're passionate about and that will resonate with your target audience. Next, invest in some basic recording equipment, such as a microphone and headphones, and find a quiet space to record. Plan your episodes in advance and create an outline to keep your content organized. Once you've recorded your episodes, use editing software to polish the audio and add any necessary music or sound effects. Finally, publish your podcast on platforms like Apple Podcasts, Spotify, and Google Podcasts, and promote it on social media to attract listeners.",
	"Growing a garden can be a rewarding and therapeutic activity, even if you have limited space. Start by choosing the right plants for your environment and the amount of sunlight your space receives. If you have a small backyard or balcony, consider using containers or vertical gardening techniques to maximize your space. Make sure to use high-quality soil and provide your plants with the necessary nutrients and water. Regularly check for pests and diseases, and take steps to prevent or treat any issues that arise. Gardening can help you relax, improve your mental health, and provide you with fresh, homegrown produce.",
	"Time management is a crucial skill for achieving your goals and maintaining a healthy work-life balance. Start by setting clear priorities and breaking your tasks into manageable chunks. Use tools like to-do lists, calendars, and time-tracking apps to stay organized and focused. Make sure to schedule regular breaks to recharge and avoid burnout. Learn to delegate tasks when possible and avoid multitasking, as it can reduce your overall productivity. By developing effective time management strategies, you can increase your efficiency, reduce stress, and achieve a better balance between work and personal life.",
	"Writing can be a powerful tool for self-expression, creativity, and personal growth. Whether you're interested in fiction, non-fiction, or poetry, there are many ways to improve your writing skills. Start by reading widely and studying the works of authors you admire. Practice writing regularly, even if it's just for a few minutes each day. Experiment with different styles and genres to find your unique voice. Join a writing group or take a writing class to receive feedback and support from other writers. Remember, the key to becoming a better writer is persistence and dedication.",
	"Sustainable living involves making choices that reduce your environmental impact and promote a healthier planet. Start by incorporating eco-friendly practices into your daily routine, such as using reusable bags, reducing water consumption, and choosing energy-efficient appliances. Consider adopting a plant-based diet, which can significantly lower your carbon footprint. Support sustainable brands and products, and educate yourself on the environmental impact of your choices. By making small changes in your lifestyle, you can contribute to a more sustainable future.",
	"The benefits of volunteering extend beyond helping others. Volunteering can provide a sense of purpose, improve your mental health, and help you develop new skills. It also offers an opportunity to connect with your community and make a positive impact. Whether you choose to volunteer at a local food bank, animal shelter, or community center, your efforts can make a difference. Consider finding volunteer opportunities that align with your interests and passions to maximize the experience.",
	"Creating a vision board can be a powerful tool for manifesting your goals and dreams. Start by gathering magazines, photos, and other materials that inspire you. Cut out images and words that resonate with your aspirations and arrange them on a board. Place your vision board in a prominent location where you'll see it daily. This constant visual reminder can help keep you motivated and focused on your goals. Remember to update your vision board regularly as your goals evolve.",
	"Photography is an art form that allows you to capture and preserve moments in time. Whether you're a beginner or an experienced photographer, there are always new techniques and skills to learn. Start by mastering the basics of composition, lighting, and exposure. Experiment with different genres of photography, such as landscape, portrait, or street photography. Practice regularly and seek feedback from other photographers to improve your skills. Remember, the key to becoming a better photographer is persistence and creativity.",
	"Self-discipline is a crucial trait for achieving your goals and maintaining a healthy lifestyle. Start by setting clear and realistic goals for yourself. Break these goals into smaller, manageable tasks and create a plan to achieve them. Hold yourself accountable by tracking your progress and celebrating your achievements. Develop healthy habits and routines that support your goals, such as regular exercise, healthy eating, and sufficient sleep. Remember, self-discipline is a skill that can be developed with practice and consistency.",
	"Exploring local cuisine is a wonderful way to experience the culture and flavors of a region. Start by researching the traditional dishes and ingredients of your destination. Visit local markets, food stalls, and restaurants to sample authentic cuisine. Don't be afraid to try new foods and flavors, even if they seem unfamiliar. Consider taking a cooking class to learn how to prepare local dishes yourself. By immersing yourself in the culinary traditions of a place, you can gain a deeper appreciation for its culture.",
	"Reading is a powerful tool for personal growth and knowledge. Make time for reading in your daily routine, even if it's just a few minutes before bed. Choose books that interest you and challenge your thinking. Join a book club or online reading group to discuss and share your thoughts with others. Keep a reading journal to track your progress and reflect on what you've learned. Remember, reading is not just about gaining knowledge but also about enjoying the journey and expanding your horizons.",
	"Fitness is an essential component of a healthy lifestyle. Incorporate regular physical activity into your routine, such as walking, running, or strength training. Find activities that you enjoy and that fit your fitness level. Set realistic fitness goals and track your progress to stay motivated. Remember to prioritize rest and recovery to prevent injuries and maintain your overall well-being. Whether you prefer working out at the gym, participating in group classes, or exercising at home, the key is to stay consistent and make fitness a priority.",
	"Mindfulness is the practice of being present and fully engaged in the moment. It can help reduce stress, improve focus, and enhance your overall well-being. Start by incorporating mindfulness exercises into your daily routine, such as deep breathing, meditation, or mindful eating. Practice paying attention to your thoughts and feelings without judgment. Use mindfulness techniques to manage stress and stay grounded in challenging situations. Remember, mindfulness is a skill that can be developed with practice and consistency.",
	"Creating a balanced diet involves making informed food choices that provide the nutrients your body needs. Focus on incorporating a variety of fruits, vegetables, whole grains, and lean proteins into your meals. Limit your intake of processed foods, added sugars, and unhealthy fats. Stay hydrated by drinking plenty of water throughout the day. Plan your meals ahead of time to ensure you're getting a balanced mix of nutrients. Remember, a healthy diet is not about restriction but about making nourishing choices that support your overall health and well-being.",
	"When it comes to survival in the wild, one of the most essential skills is the ability to build a shelter. A good shelter protects you from the elements, keeps you warm, and provides a sense of security. Start by finding a suitable location, preferably on high ground and away from water sources to avoid flooding. Look for natural materials such as branches, leaves, and grass to construct your shelter. There are various types of shelters you can build, including lean-tos, debris huts, and A-frame shelters. Each type has its advantages and can be adapted based on the available materials and weather conditions. Remember to insulate the ground and walls to retain heat and keep you comfortable throughout the night.",
	"Finding and purifying water is crucial for survival. In the wild, natural water sources such as rivers, lakes, and streams can be contaminated with harmful bacteria and parasites. To ensure the water is safe to drink, you need to purify it. Boiling is one of the most effective methods, as it kills most pathogens. If you don't have a pot, you can use a metal container or even a hollowed-out log to hold the water over a fire. Other purification methods include using water purification tablets, portable water filters, or UV light devices. In emergency situations, you can also use improvised methods such as solar stills or collecting rainwater. Always remember that staying hydrated is vital for your survival, so make finding clean water a top priority.",
	"Knowing how to start a fire is a fundamental survival skill. A fire provides warmth, cooks food, purifies water, and can signal for help. To start a fire, you'll need tinder, kindling, and fuel. Tinder consists of small, dry materials that ignite easily, such as dry grass, leaves, or bark. Kindling is slightly larger, dry sticks and twigs that catch fire from the tinder. Fuel includes larger pieces of wood that sustain the fire. There are various methods to start a fire, including using a fire starter, flint and steel, or a bow drill. Each method requires practice and patience, so it's essential to practice fire-starting techniques before you find yourself in a survival situation. Remember to build your fire in a safe location, away from flammable materials, and always extinguish it thoroughly before leaving.",
	"Foraging for wild edibles is a valuable skill that can provide you with essential nutrients in a survival situation. Before foraging, it's crucial to familiarize yourself with the local flora and learn to identify edible plants, fruits, nuts, and mushrooms. Always follow the rule of foraging: if in doubt, leave it out. Some plants have poisonous look-alikes, so it's essential to be certain of the identification before consuming any wild food. It's also important to forage sustainably, taking only what you need and leaving enough for the ecosystem to thrive. Some common wild edibles include dandelions, wild garlic, cattails, and berries. Incorporating foraged food into your diet can boost your nutrition and increase your chances of survival in the wild.",
	"Navigation skills are essential for finding your way in the wilderness. Without a map and compass, you can still use natural navigation techniques to stay oriented. The sun, stars, and natural landmarks can all serve as navigational aids. During the day, the sun rises in the east and sets in the west, providing a general direction. At night, the North Star (Polaris) can help you find true north. Natural landmarks such as mountains, rivers, and unique rock formations can also guide you. Learning to read the landscape and observe your surroundings will help you stay on course. If you do have a map and compass, practice using them to navigate accurately. Always plan your route, set checkpoints, and stay aware of your location to avoid getting lost.",
	"First aid knowledge is critical in survival situations. Accidents and injuries can happen, and knowing how to treat them can save lives. Basic first aid skills include treating cuts and wounds, setting broken bones, and performing CPR. Always carry a well-stocked first aid kit with bandages, antiseptics, pain relievers, and any necessary medications. Learn how to identify and treat common wilderness ailments such as hypothermia, heatstroke, and dehydration. In a survival scenario, improvisation may be necessary, using natural materials to create splints or dressings. Staying calm and assessing the situation is key to providing effective first aid. Remember, your ability to handle medical emergencies can significantly impact your chances of survival.",
	"Hunting and fishing are essential skills for sourcing food in the wild. Small game hunting can provide a steady supply of protein. Learn to set traps and snares for animals such as rabbits, squirrels, and birds. Building a survival bow and arrows or using a slingshot can also be effective hunting methods. Fishing can be done with minimal equipment, using a simple fishing line, hooks, and bait. Look for bodies of water with fish activity, such as lakes, rivers, and streams. Insects, worms, and small fish can serve as bait. Patience is crucial, as successful hunting and fishing require time and persistence. Understanding animal behavior and tracking skills can also increase your chances of capturing food.",
	"Signaling for help is a crucial aspect of survival, especially if you're lost or injured. There are various methods to signal for rescue, including using fire, mirrors, whistles, and signal flags. Building a large, smoky fire can attract attention from a distance. Mirrors or reflective surfaces can be used to flash sunlight toward potential rescuers. Three short blasts from a whistle or three signal fires arranged in a triangle are recognized as distress signals. Brightly colored clothing or signal flags can also catch the attention of passing aircraft or search parties. Always remain in open areas where you can be easily spotted and stay calm while waiting for help to arrive. Knowing how to effectively signal for help can significantly increase your chances of being rescued.",
	"In a survival situation, maintaining a positive mindset is just as important as physical skills. Staying calm and focused can help you make rational decisions and avoid panic. Practice stress-reduction techniques such as deep breathing, visualization, and mindfulness. Set small, achievable goals to keep yourself motivated and maintain a sense of control. Staying positive can also help you conserve energy and think creatively. Remember that survival is often a mental challenge as much as a physical one. Keeping a positive outlook and believing in your ability to overcome obstacles can increase your chances of success in a survival scenario.",
	"Building a proper survival kit is essential for preparedness. Your kit should include items such as a multi-tool, fire starter, water purification tablets, first aid supplies, and emergency food rations. A lightweight, durable backpack can store all your essentials. Customize your kit based on the environment you'll be in, such as including insect repellent for the tropics or extra thermal layers for cold climates. Regularly check and update your kit to ensure all items are in good condition and replace any expired supplies. A well-prepared survival kit can make a significant difference in your ability to navigate and endure a challenging situation.",
	"Cooking is an art that brings people together and creates memorable experiences. One of the most important aspects of cooking is understanding the basics, such as knife skills, cooking techniques, and the importance of fresh ingredients. Start by investing in quality kitchen tools and equipment. Learn how to properly use a chef's knife to chop, slice, and dice with precision. Master basic cooking techniques like sautéing, roasting, and boiling. Understanding these fundamentals will give you the confidence to experiment with new recipes and flavors. Don't be afraid to make mistakes, as they are an essential part of the learning process. Embrace the joy of cooking and the satisfaction of creating delicious meals from scratch.",
	"Meal planning is a great way to save time and ensure you always have a nutritious meal on hand. Start by creating a weekly meal plan that includes a variety of dishes to keep things interesting. Consider your dietary needs and preferences, and aim for a balanced mix of proteins, vegetables, and whole grains. Make a shopping list based on your meal plan to avoid unnecessary trips to the store. Spend a few hours each week prepping ingredients and cooking meals in advance. This can include chopping vegetables, marinating proteins, and cooking grains. Store your prepped meals in airtight containers to keep them fresh throughout the week. Meal planning not only saves time but also helps you make healthier food choices and reduces food waste.",
	"One of the most rewarding aspects of cooking is experimenting with different cuisines from around the world. Each culture has its unique flavors, ingredients, and cooking techniques that can broaden your culinary horizons. Start by exploring popular dishes from different regions, such as Italian pasta, Indian curries, or Japanese sushi. Learn about the key ingredients and spices used in each cuisine, and don't be afraid to try new cooking methods. Experimenting with international flavors can inspire creativity in the kitchen and help you discover new favorite dishes. Whether you're cooking a Moroccan tagine or a Thai stir-fry, embracing global cuisines can bring excitement and diversity to your meals.",
	"Baking is a precise science that requires attention to detail and patience. Understanding the chemistry behind baking can help you achieve perfect results every time. Start by familiarizing yourself with common baking ingredients, such as flour, sugar, eggs, and leavening agents. Learn the roles each ingredient plays in the baking process and how they interact with one another. Follow recipes closely and measure ingredients accurately to ensure consistency. Experiment with different baking techniques, such as creaming, folding, and kneading, to achieve the desired texture and flavor. Whether you're baking bread, cookies, or cakes, mastering the basics of baking can open up a world of delicious possibilities.",
	"Healthy cooking doesn't have to be bland or boring. By making simple swaps and incorporating nutrient-dense ingredients, you can create delicious and wholesome meals. Start by focusing on whole foods, such as fresh fruits, vegetables, lean proteins, and whole grains. Avoid processed foods and added sugars, which can detract from the nutritional value of your meals. Experiment with different cooking methods, such as grilling, steaming, and roasting, to enhance the flavors of your ingredients. Use herbs and spices to add depth and complexity to your dishes without relying on excess salt or fat. By making mindful choices in the kitchen, you can enjoy a healthy and flavorful diet that supports your overall well-being.",
	"Cooking for special occasions and holidays can be both exciting and challenging. Whether you're hosting a dinner party or preparing a festive feast, planning and organization are key to success. Start by creating a menu that includes a variety of dishes to accommodate different tastes and dietary needs. Consider incorporating seasonal ingredients to make your dishes feel fresh and timely. Make a detailed shopping list and gather all the necessary ingredients in advance. Plan your cooking schedule to ensure that everything is ready on time. Don't be afraid to ask for help from friends and family to lighten the workload. Remember, the goal is to create a memorable and enjoyable experience for your guests, so focus on the joy of sharing good food and company.",
	"Cooking with kids can be a fun and educational experience that teaches valuable life skills. Involving children in the kitchen encourages creativity, responsibility, and a love for healthy eating. Start by choosing age-appropriate tasks for your kids, such as washing vegetables, stirring ingredients, or measuring out portions. Teach them basic kitchen safety rules, such as handling knives carefully and using oven mitts. Encourage them to explore different flavors and textures by allowing them to taste and smell ingredients. Cooking together can also be a great opportunity to introduce new foods and discuss the importance of a balanced diet. By making cooking a family activity, you can create lasting memories and instill a lifelong appreciation for good food.",
	"Understanding food safety is crucial for preventing foodborne illnesses and ensuring the health and well-being of your family. Start by practicing proper hygiene, such as washing your hands thoroughly before and after handling food. Keep your kitchen clean and sanitized, and regularly clean cutting boards, utensils, and countertops. Learn about the safe handling and storage of different types of food, such as raw meat, poultry, and seafood. Use separate cutting boards for raw and cooked foods to prevent cross-contamination. Cook foods to the appropriate internal temperatures to kill harmful bacteria. By following these food safety guidelines, you can enjoy cooking with confidence and protect your loved ones from potential hazards.",
	"One of the most enjoyable aspects of cooking is sharing your creations with others. Hosting a dinner party or potluck can be a great way to connect with friends and family over good food. Start by planning a menu that includes a variety of dishes to suit different tastes and dietary needs. Consider asking your guests to contribute a dish, which can add diversity to the meal and lighten your workload. Set the table with care, using nice dishes, cutlery, and decorations to create a welcoming atmosphere. Encourage conversation and interaction by serving dishes family-style or creating a buffet. Remember, the goal is to create a warm and inviting environment where everyone can enjoy good food and good company.",
	"Cooking can be a therapeutic and stress-relieving activity that allows you to unwind and focus on the present moment. The act of chopping vegetables, stirring a pot, or kneading dough can be meditative and calming. To enhance the therapeutic benefits of cooking, create a pleasant and organized kitchen environment. Play your favorite music, light a scented candle, or open a window to let in fresh air. Take your time with each step of the cooking process, and enjoy the sensory experiences of touch, smell, and taste. Cooking mindfully can help you relax, reduce stress, and cultivate a sense of accomplishment and satisfaction. Whether you're preparing a simple meal or an elaborate dish, embrace the joy and therapeutic benefits of cooking.",
	"In the ever-evolving world of software development, staying updated with the latest trends and technologies is crucial. One such trend that has gained significant traction in recent years is DevOps. DevOps is a set of practices that combines software development (Dev) and IT operations (Ops) with the goal of shortening the development lifecycle and providing continuous delivery of high-quality software. By fostering a culture of collaboration and automation, DevOps helps teams to improve efficiency, reduce errors, and deploy features faster. Key practices include continuous integration (CI), continuous delivery (CD), automated testing, and infrastructure as code (IaC). Embracing DevOps can lead to more reliable software releases and increased customer satisfaction.",
	"Microservices architecture has become a popular approach for building scalable and maintainable software systems. Unlike traditional monolithic architecture, where all components are tightly coupled, microservices break down an application into smaller, loosely coupled services. Each service is responsible for a specific functionality and can be developed, deployed, and scaled independently. This approach offers several benefits, including improved flexibility, faster development cycles, and easier maintenance. However, it also introduces challenges such as increased complexity in managing inter-service communication and ensuring data consistency. To successfully implement microservices, it's essential to adopt best practices such as proper service design, robust API management, and effective monitoring and logging.",
	"Continuous integration and continuous delivery (CI/CD) are foundational practices in modern software development. CI involves automatically integrating code changes from multiple contributors into a shared repository several times a day. This practice helps to identify and fix integration issues early, reducing the risk of conflicts and bugs. CD takes CI a step further by automating the deployment process, ensuring that code changes are consistently and reliably delivered to production. Implementing a CI/CD pipeline typically involves using tools like Jenkins, GitLab CI, or CircleCI, along with version control systems such as Git. By adopting CI/CD, development teams can achieve faster release cycles, higher code quality, and greater confidence in their software.",
	"In recent years, the rise of artificial intelligence (AI) and machine learning (ML) has significantly impacted the software development landscape. AI and ML enable developers to create applications that can learn from data, make predictions, and improve over time. Common use cases include recommendation systems, natural language processing, and image recognition. To get started with AI and ML, developers need to understand the fundamentals of data science, algorithms, and statistical modeling. Popular frameworks and libraries like TensorFlow, PyTorch, and Scikit-learn provide the tools necessary to build and train machine learning models. Integrating AI and ML into your software can open up new possibilities and create more intelligent and adaptive applications.",
	"Test-driven development (TDD) is a software development approach that emphasizes writing tests before writing the actual code. The TDD process involves writing a test for a specific functionality, running the test to ensure it fails (since the functionality isn't implemented yet), writing the minimum amount of code to make the test pass, and then refactoring the code for improvement. This cycle is repeated for each new feature or functionality. TDD helps to ensure that code is thoroughly tested and encourages better design and modularity. It also provides a safety net for refactoring and helps developers catch bugs early in the development process. While TDD requires discipline and practice, it can lead to higher code quality and more maintainable software.",
	"Version control is a critical aspect of software development, enabling teams to track changes, collaborate effectively, and manage codebases efficiently. Git is one of the most widely used version control systems, known for its speed, flexibility, and distributed nature. With Git, developers can work on branches to isolate their changes, merge changes from different branches, and resolve conflicts. Key commands include git init, git clone, git add, git commit, git push, and git pull. Platforms like GitHub, GitLab, and Bitbucket provide additional features for collaboration, such as pull requests, code reviews, and issue tracking. Understanding and mastering version control is essential for any software developer to manage code changes and collaborate seamlessly.",
	"Agile methodologies have transformed the way software development teams work, promoting flexibility, collaboration, and iterative progress. The Agile Manifesto, published in 2001, outlines values and principles that prioritize individuals and interactions, working software, customer collaboration, and responding to change. Popular Agile frameworks include Scrum, Kanban, and Extreme Programming (XP). Scrum involves working in time-boxed iterations called sprints, with roles such as Scrum Master and Product Owner. Kanban focuses on visualizing work and limiting work in progress (WIP) to improve flow. XP emphasizes technical practices like pair programming, continuous integration, and test-driven development. By adopting Agile methodologies, teams can deliver value more frequently, adapt to changing requirements, and improve overall efficiency.",
	"Code reviews are an essential practice in software development, helping to maintain code quality, share knowledge, and catch bugs early. During a code review, developers examine each other's code changes, providing feedback and suggesting improvements. Effective code reviews involve clear communication, constructive feedback, and a focus on both the functional and non-functional aspects of the code. Tools like GitHub, GitLab, and Bitbucket offer integrated code review features, allowing developers to comment on specific lines of code, request changes, and approve pull requests. By incorporating regular code reviews into the development process, teams can improve code quality, foster collaboration, and ensure that best practices are followed.",
	"The rise of containerization has revolutionized the way software is developed, deployed, and managed. Containers are lightweight, portable units that package an application and its dependencies, ensuring consistency across different environments. Docker is one of the most popular containerization platforms, allowing developers to create, deploy, and run applications in containers. Kubernetes, an open-source container orchestration platform, simplifies the management of containerized applications, automating tasks such as scaling, load balancing, and self-healing. By adopting containerization and orchestration tools, development teams can achieve greater flexibility, scalability, and efficiency in their software delivery processes.",
	"Security is a paramount concern in software development, as vulnerabilities can lead to data breaches, financial losses, and reputational damage. Secure coding practices help to minimize risks and protect applications from common threats such as SQL injection, cross-site scripting (XSS), and buffer overflow attacks. Key practices include input validation, proper authentication and authorization mechanisms, secure data storage, and regular security testing. Tools like static code analyzers, vulnerability scanners, and penetration testing frameworks can help identify and remediate security issues. By prioritizing security throughout the development lifecycle, teams can build more robust and resilient software that protects users and their data.",
}

var tags = []string{
	"travel", "food", "recipe", "technology", "lifestyle", "health", "fitness", "fashion", "beauty", "finance",
	"business", "education", "marketing", "SEO", "socialmedia", "photography", "art", "design", "writing", "books",
	"parenting", "pets", "DIY", "home", "garden", "decor", "realestate", "automotive", "sports", "music",
	"movies", "TV", "entertainment", "gaming", "science", "nature", "environment", "politics", "news", "history",
	"culture", "travelguide", "travelhacks", "foodie", "recipeideas", "technews", "lifestyleblog", "healthtips", "workout", "style",
	"makeup", "skincare", "money", "investing", "smallbusiness", "startup", "educationtips", "digitalmarketing", "SEOtips", "socialmediamarketing",
	"phototips", "artwork", "graphicdesign", "writingtips", "bookreview", "parentingadvice", "petcare", "DIYprojects", "homeimprovement", "gardening",
	"homedecor", "realestatetips", "carreviews", "sportsnews", "musictips", "moviereview", "TVshows", "entertainmentnews", "gamingtips", "sciencenews",
	"wildlife", "conservation", "politicsnews", "worldnews", "historyfacts", "culturetips", "travelphotography", "foodblog", "recipeblog", "techblog",
	"lifestyleblogger", "healthblog", "fitnessblog", "fashionblog", "beautyblog", "financeblog", "businessblog", "educationblog", "marketingblog", "SEOblog",
	"socialmediablog", "photographyblog", "artblog", "designblog", "writingblog", "bookblog", "parentingblog", "petsblog", "DIYblog", "homeblog",
	"gardenblog", "decorblog", "realestateblog", "automotiveblog", "sportsblog", "musicblog", "movieblog", "TVblog", "entertainmentblog", "gamingblog",
	"scienceblog", "natureblog", "environmentblog", "politicsblog", "newsblog", "historyblog", "cultureblog", "travelstories", "foodstories", "recipestories",
	"techstories", "lifestylestories", "healthstories", "fitnessstories", "fashionstories", "beautystories", "financestories", "businessstories", "educationstories", "marketingstories",
	"SEOstories", "socialmediastories", "photographystories", "artstories", "designstories", "writingstories", "bookstories", "parentingstories", "petstories", "DIYstories",
	"homestories", "gardenstories", "decorstories", "realestatestories", "automotivestories", "sportsstories", "musicstories", "moviestories", "TVstories", "entertainmentstories",
	"gamingstories", "sciencestories", "naturestories", "environmentstories", "politicsstories", "newsstories", "historystories", "culturestories", "traveltips", "foodtips",
	"recipetips", "techtips", "lifestyletips", "healthtips", "fitnesstips", "fashiontips", "beautytips", "financetips", "businesstips", "educationtips",
	"marketingtips", "SEOtips", "socialmediatips", "photographytips", "arttips", "designtips", "writingtips", "booktips", "parentingtips", "pettips",
	"DIYtips", "hometips", "gardentips", "decortips", "realestatetips", "automotivetips", "sportstips", "musictips", "movietips", "TVtips",
	"entertainmenttips", "gamingtips", "sciencetips", "naturetips", "environmenttips", "politicstips", "newstips", "historytips", "culturetips",
}

var comments = []string{
	"Fantastic read! Thank you.",
	"This article is very informative.",
	"Really enjoyed this post!",
	"Great insights, thanks for sharing.",
	"This was incredibly helpful.",
	"Well-explained, thank you!",
	"I appreciate this information.",
	"Thanks for the tips!",
	"Excellent article!",
	"Very useful content.",
	"I learned something new today.",
	"I completely agree with this.",
	"Wonderful post!",
	"Great job on this article.",
	"Thanks for the advice!",
	"Very well written.",
	"I found this very insightful.",
	"This was a great read!",
	"Thanks for sharing your thoughts.",
	"I appreciate the clarity.",
	"This post is spot on.",
	"Thanks for the detailed explanation.",
	"I love this topic.",
	"Great tips, very useful.",
	"This was exactly what I needed.",
	"Very informative article.",
	"Great write-up!",
	"This made my day.",
	"Very practical advice.",
	"Thanks for the recommendations.",
	"This is very helpful.",
	"Nicely done!",
	"Really useful information.",
	"I appreciate your perspective.",
	"Great read, thanks!",
	"Very well articulated.",
	"I found this really useful.",
	"This is a must-read.",
	"Thanks for breaking it down.",
	"Fantastic insights.",
	"This is so timely.",
	"Thanks for sharing!",
	"Love this!",
	"Great job!",
	"This was very enlightening.",
	"I appreciate this content.",
	"Wonderful write-up.",
	"This was very interesting.",
	"Excellent tips!",
	"Very practical.",
	"Great advice, thanks!",
	"I found this very helpful.",
	"Thanks for the clear explanation.",
	"This is a fantastic resource.",
	"Very useful advice.",
	"I appreciate the thoroughness.",
	"Thanks for the valuable info.",
	"This is very well explained.",
	"Great job with this post.",
	"Thanks for the insights!",
	"I love this post.",
	"Very detailed and informative.",
	"This was very useful.",
	"Great information, thanks!",
	"I learned a lot from this.",
	"Thanks for sharing your expertise.",
	"This is a great article.",
	"Very helpful, thank you!",
	"I appreciate the effort.",
	"Thanks for the awesome tips.",
	"This was very informative.",
	"Great content!",
	"This post is very useful.",
	"I really enjoyed this read.",
	"Thanks for the advice.",
	"Very clear and concise.",
	"I found this very insightful.",
	"This was a fantastic read.",
	"Thanks for the tips and tricks.",
	"Great explanation!",
	"This post is a gem.",
	"Very practical information.",
	"Thanks for the helpful advice.",
	"I love the insights here.",
	"Well-written and informative.",
	"This was incredibly useful.",
	"Thanks for the great post!",
	"Very thoughtful article.",
	"I appreciate the detail.",
	"This is very well written.",
	"Great read, thank you!",
	"Thanks for the practical tips.",
	"This was very enlightening.",
	"I found this very useful.",
	"Great job with this article.",
	"Very informative and clear.",
	"I love this content.",
	"This was very helpful.",
	"Great advice!",
	"Thanks for the guidance.",
	"This post is very timely.",
	"Very helpful, thanks!",
	"I learned a lot.",
	"Fantastic job on this post.",
	"Very insightful content.",
	"I appreciate the thorough explanation.",
	"Thanks for the tips!",
	"Great post, very useful.",
	"I found this very practical.",
	"Thanks for sharing your knowledge.",
	"This was a great read.",
	"Very well explained.",
	"I appreciate this article.",
	"Great advice, thank you!",
	"This is very useful information.",
	"Thanks for the helpful post.",
	"I learned something new.",
	"Great job explaining this.",
	"Very detailed and clear.",
	"This was very informative.",
	"Thanks for the tips and advice.",
	"Great insights!",
	"This was very helpful.",
	"I love this article.",
	"Thanks for the practical information.",
	"Very useful content.",
	"I appreciate this advice.",
	"Great post, thanks for sharing.",
	"This was a fantastic read.",
	"I found this very enlightening.",
	"Thanks for the detailed guide.",
	"Great job with this article.",
	"This is very well written.",
	"Thanks for the helpful tips.",
	"I love this content.",
	"Very informative post.",
	"This was very useful.",
	"Thanks for the clear explanation.",
	"Great advice, very helpful.",
	"I appreciate the insights.",
	"This was very enlightening.",
	"I found this very practical.",
	"Thanks for sharing this post.",
	"Great job!",
	"This was very informative.",
	"I love the tips here.",
	"Very useful advice.",
	"Thanks for the great post!",
	"I appreciate the thoroughness.",
	"This was very helpful.",
	"Great insights, thank you!",
	"This is very well explained.",
	"Thanks for the advice.",
	"Very useful post.",
	"I appreciate this content.",
	"Great article!",
	"This was very helpful.",
	"Thanks for the tips and tricks.",
	"Great job explaining this.",
	"Very informative and clear.",
	"I love this post.",
	"This was a fantastic read.",
	"Thanks for the practical advice.",
	"Great insights!",
	"This is very well written.",
	"Thanks for the guidance.",
	"I found this very helpful.",
	"Very detailed article.",
	"Great post, thank you!",
	"I appreciate the explanation.",
	"Thanks for sharing this information.",
	"This was very useful.",
	"Great advice, thank you!",
	"This post is a gem.",
	"Very practical tips.",
	"I appreciate this article.",
	"Thanks for the detailed explanation.",
	"Great job with this post.",
	"This was very helpful.",
	"I love the insights here.",
	"Thanks for the practical information.",
	"Very useful content.",
	"I appreciate this advice.",
	"Thanks for the clear explanation.",
	"This post is very useful.",
	"I learned a lot from this.",
	"Great job explaining this.",
	"Very informative post.",
	"I love this article.",
}

func Seed(repo repo.Repository) {

	// generate all users - they referenced used randomly below
	users := generateUsers(len(userNames))
	for _, user := range users {
		log.Println(user.Username)
		if err := repo.UserSeed.Create(user); err != nil {
			log.Println("error seeding users:", err)
			return
		}
	}

	posts := generatePosts(100, users)
	for _, post := range posts {
		log.Println(fmt.Sprintf("user: %d post: %s", post.UserID, post.Title))
		if err := repo.Posts.Create(post); err != nil {
			log.Println("error seeding posts", err)
			return
		}
	}

	cms := generateComments(300, users, posts)
	for _, comment := range cms {
		log.Println(fmt.Sprintf("user: %d comment: %s", comment.UserID, comment.Content))
		if err := repo.Comments.Create(comment); err != nil {
			log.Println("error seeding comments", err)
			return
		}
	}

	log.Println("Seed completed")
}

func generateUsers(num int) []*repo.User {
	users := make([]*repo.User, num)

	for i := 0; i < num; i++ {
		//userName := userNames[i%len(userNames)] + fmt.Sprintf("%d", i)
		userName := userNames[i]
		users[i] = &repo.User{
			Username: userName,
			Email:    userName + "@" + domainNames[rand.Intn(len(domainNames))],
			RoleID:   1,
			//Password: "password",
		}
	}

	return users
}

func generatePosts(num int, users []*repo.User) []*repo.Post {
	posts := make([]*repo.Post, num)

	for i := 0; i < num; i++ {
		user := users[rand.Intn(len(users))]
		posts[i] = &repo.Post{
			UserID:  user.ID,
			Title:   titles[rand.Intn(len(titles))],
			Content: contents[rand.Intn(len(contents))],
			Tags:    tags[rand.Intn(len(tags))] + ", " + tags[rand.Intn(len(tags))],
		}
	}

	return posts
}

func generateComments(num int, users []*repo.User, posts []*repo.Post) []*repo.Comment {
	cms := make([]*repo.Comment, num)

	for i := 0; i < num; i++ {
		cms[i] = &repo.Comment{
			UserID:  users[rand.Intn(len(users))].ID,
			PostID:  posts[rand.Intn(len(posts))].ID,
			Content: comments[rand.Intn(len(comments))],
		}
	}

	return cms
}
