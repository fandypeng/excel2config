package helper

import (
	"go.mongodb.org/mongo-driver/x/mongo/driver/uuid"
	"math/rand"
)

var avatarUrlList = []string{
	"http://wespynextpic.afunapp.com/o_1d8iju3n19e1hhu1vgvl4q92b1s.jpg",
	"http://wespynextpic.afunapp.com/o_1d8iju3n17c2mne1uft1al11kva1t.jpg",
	"http://wespynextpic.afunapp.com/o_1d8iju3n2c5m1okj1hid5dm1j8u1u.jpg",
	"http://wespynextpic.afunapp.com/o_1d8iju3n21vtc3u2918188u1chi1v.jpg",
	"http://wespynextpic.afunapp.com/o_1d8iju3n2d2q72n1c93ivn13j120.jpg",
	"http://wespynextpic.afunapp.com/o_1d8iju3n212fr1um7klr1j9do8921.jpg",
	"http://wespynextpic.afunapp.com/o_1d8iju3n2pk26t18b131n19t422.jpg",
	"http://wespynextpic.afunapp.com/o_1d8iju3n21t291cdj1e301loq15ep23.jpg",
	"http://wespynextpic.afunapp.com/o_1d8iju3n2ve116ug8km46h1tai24.jpg",
	"http://wespynextpic.afunapp.com/o_1d8iju3n23om1g8728ubgden25.jpg",
	"http://wespynextpic.afunapp.com/o_1d8iju3n2c0l1fu1pktebu1vmf26.jpg",
	"http://wespynextpic.afunapp.com/o_1d8iju3n21gtf1m4d12ndilq1q3527.jpg",
	"http://wespynextpic.afunapp.com/o_1d8iju3n2cksc501bdq171j1h5m28.jpg",
	"http://wespynextpic.afunapp.com/o_1d8iju3n21hptv1h1vc4ro0e5b29.jpg",
	"http://wespynextpic.afunapp.com/o_1d8iju3n22l7ilt1e6t17ae1uvk2a.jpg",
	"http://wespynextpic.afunapp.com/o_1d8iju3n24011rv9rb3j2o1if62b.jpg",
	"http://wespynextpic.afunapp.com/o_1d8iju3n21i4j1aot1uu95j713032c.jpg",
	"http://wespynextpic.afunapp.com/o_1d8iju3n2nla2pp1kd5mbh1fof2d.jpg",
	"http://wespynextpic.afunapp.com/o_1d8iju3n3f00u1h1cn2o6017lp2e.jpg",
	"http://wespynextpic.afunapp.com/o_1d8iju3n31qp9rvu1raq1vq216u72f.jpg",
	"http://wespynextpic.afunapp.com/o_1d8iju3n3ofn1gh6alrest1urs2g.jpg",
	"http://wespynextpic.afunapp.com/o_1d8iju3n3cja1coch5g321ome2h.jpg",
	"http://wespynextpic.afunapp.com/o_1d8iju3n3ns1a201c1uesuer72i.jpg",
	"http://wespynextpic.afunapp.com/o_1d8iju3n3c1a1p791a06v753ue2j.jpg",
	"http://wespynextpic.afunapp.com/o_1d8iju3n312uv119hd88qv1146j2k.jpg",
	"http://wespynextpic.afunapp.com/o_1d8iju3n31kqe1f9kvsk1kf51k4g2l.jpg",
	"http://wespynextpic.afunapp.com/o_1d8iju3n38ms13m013n917nth4c2m.jpg",
	"http://wespynextpic.afunapp.com/o_1d8iju3n31stsovl2ufqcspjn2n.jpg",
	"http://wespynextpic.afunapp.com/o_1d8iju3n3d7bagplv610ib3772o.jpg",
	"http://wespynextpic.afunapp.com/o_1d8iju3n3q311tgqjjq1b2a1i8i2p.jpg",
	"http://wespynextpic.afunapp.com/o_1d8iju3n346s105v1it7vvm1c2f2q.jpg",
	"http://wespynextpic.afunapp.com/o_1d8iju3n34tl1fita9a1tj25tu2r.jpg",
	"http://wespynextpic.afunapp.com/o_1d8iju3n31k051q632cepp0toh2s.jpg",
	"http://wespynextpic.afunapp.com/o_1d8iju3n31q94igi4fni7o7372t.jpg",
	"http://wespynextpic.afunapp.com/o_1d8iju3n3f1r2hg1fjb1pvov702u.jpg",
	"http://wespynextpic.afunapp.com/o_1d8iju3n319kn1mq1vhmkh118q2v.jpg",
	"http://wespynextpic.afunapp.com/o_1d8iju3n31nb311sh83sef13j530.jpg",
	"http://wespynextpic.afunapp.com/o_1d8iju3n3he11pp61e4a1k7214t231.jpg",
	"http://wespynextpic.afunapp.com/o_1d8iju3n3kd71nko6gh1dag1pmk32.jpg",
	"http://wespynextpic.afunapp.com/o_1d8iju3n3ptn13cs1dbdfp31j5833.jpg",
	"http://wespynextpic.afunapp.com/o_1d8iju3n3f351atdtkr1j5ag2234.jpg",
	"http://wespynextpic.afunapp.com/o_1d8iju3n311091t2trhg1c64c7s35.jpg",
	"http://wespynextpic.afunapp.com/o_1d8iju3n315qdgro180tvjf1tid36.jpg",
	"http://wespynextpic.afunapp.com/o_1d8iju3n318r319911th17qh16gk37.jpg",
	"http://wespynextpic.afunapp.com/o_1d8iju3n31uj11k7lgvotvc1pc338.jpg",
	"http://wespynextpic.afunapp.com/o_1d8iju3n31gfl1rv91ijni1deh39.jpg",
	"http://wespynextpic.afunapp.com/o_1d8iju3n3h1epf23rh1tbo157n3a.jpg",
	"http://wespynextpic.afunapp.com/o_1d8iju3n3ckc1d4s16ughm11a9b3b.jpg",
	"http://wespynextpic.afunapp.com/o_1d8iju3n319h91il7oih1jdg1s1n3c.jpg",
	"http://wespynextpic.afunapp.com/o_1d8iju3n3at3qv293511a27ur3d.jpg",
	"http://wespynextpic.afunapp.com/o_1d8iju3n379n1mppbgj18gj145h3e.jpg",
}

func GetRandomAvatar() string {
	if len(avatarUrlList) == 0 {
		return ""
	}
	randId := rand.Intn(len(avatarUrlList))
	return avatarUrlList[randId]
}

func GetRandomToken(uid string) string {
	u, _ := uuid.New()
	r := uid + string(u[:])
	return Md5Sum(r)
}
