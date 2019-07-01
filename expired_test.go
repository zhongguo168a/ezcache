package ezcache

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

func TestCacheSync_Expire(t *testing.T) {
	Convey("输入了正确的参数，成功返回", t, func() {
		cache := NewCacheSyncNoExpire()
		cache.Set("testkey-1", 0)
		cache.Set("testkey-2", 0)
		cache.Set("testkey-3", 0)
		cache.Expire("testkey-1", time.Millisecond*100)
		cache.Expire("testkey-2", time.Millisecond*500)
		cache.Expire("testkey-3", time.Millisecond*300)

		So(len(cache.expiredSet), ShouldEqual, 3)
		So(cache.expiredRoot.key == "testkey-1", ShouldBeTrue)

		So(cache.expiredSet["testkey-1"].prev == nil, ShouldBeTrue)
		So(cache.expiredSet["testkey-1"].next != nil && cache.expiredSet["testkey-1"].next.key == "testkey-3", ShouldBeTrue)

		So(cache.expiredSet["testkey-3"].prev != nil && cache.expiredSet["testkey-3"].prev.key == "testkey-1", ShouldBeTrue)
		So(cache.expiredSet["testkey-3"].next != nil && cache.expiredSet["testkey-3"].next.key == "testkey-2", ShouldBeTrue)

		So(cache.expiredSet["testkey-2"].prev != nil && cache.expiredSet["testkey-2"].prev.key == "testkey-3", ShouldBeTrue)
		So(cache.expiredSet["testkey-2"].next == nil, ShouldBeTrue)

		//
		time.Sleep(time.Millisecond * 200)
		cache.cleanExpireList()
		So(len(cache.expiredSet), ShouldEqual, 2)
		So(cache.expiredRoot.key == "testkey-3", ShouldBeTrue)

		So(cache.expiredSet["testkey-3"].prev == nil, ShouldBeTrue)
		So(cache.expiredSet["testkey-3"].next != nil && cache.expiredSet["testkey-3"].next.key == "testkey-2", ShouldBeTrue)

		So(cache.expiredSet["testkey-2"].prev != nil && cache.expiredSet["testkey-2"].prev.key == "testkey-3", ShouldBeTrue)
		So(cache.expiredSet["testkey-2"].next == nil, ShouldBeTrue)

		//
		time.Sleep(time.Millisecond * 200)
		cache.cleanExpireList()
		So(len(cache.expiredSet), ShouldEqual, 1)
		So(cache.expiredRoot.key == "testkey-2", ShouldBeTrue)

		So(cache.expiredSet["testkey-2"].prev == nil, ShouldBeTrue)
		So(cache.expiredSet["testkey-2"].next == nil, ShouldBeTrue)

		//
		time.Sleep(time.Millisecond * 200)
		cache.cleanExpireList()
		So(len(cache.expiredSet), ShouldEqual, 0)
		So(cache.expiredRoot == nil, ShouldBeTrue)

	})

}

func TestCacheSync_Expire2(t *testing.T) {
	Convey("输入了正确的参数，成功返回", t, func() {
		cache := NewCacheSyncNoExpire()
		cache.Set("testkey-1", 0)
		cache.Set("testkey-2", 0)
		cache.Expire("testkey-1", time.Millisecond*300)
		cache.Expire("testkey-2", time.Millisecond*100)

		So(len(cache.expiredSet), ShouldEqual, 2)
		So(cache.expiredRoot.key == "testkey-2", ShouldBeTrue)

		So(cache.expiredSet["testkey-2"].prev == nil, ShouldBeTrue)
		So(cache.expiredSet["testkey-2"].next != nil && cache.expiredSet["testkey-2"].next.key == "testkey-1", ShouldBeTrue)

		So(cache.expiredSet["testkey-1"].prev != nil && cache.expiredSet["testkey-1"].prev.key == "testkey-2", ShouldBeTrue)
		So(cache.expiredSet["testkey-1"].next == nil, ShouldBeTrue)

		//
		cache.Expire("testkey-2", time.Millisecond*500)

		So(len(cache.expiredSet), ShouldEqual, 2)
		So(cache.expiredRoot.key == "testkey-1", ShouldBeTrue)

		So(cache.expiredSet["testkey-1"].prev == nil, ShouldBeTrue)
		So(cache.expiredSet["testkey-1"].next != nil && cache.expiredSet["testkey-1"].next.key == "testkey-2", ShouldBeTrue)

		So(cache.expiredSet["testkey-2"].prev != nil && cache.expiredSet["testkey-2"].prev.key == "testkey-1", ShouldBeTrue)
		So(cache.expiredSet["testkey-2"].next == nil, ShouldBeTrue)
	})

}
