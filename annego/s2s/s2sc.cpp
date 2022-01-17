#define _GLIBCXX_USE_CXX11_ABI 0
#include "s2sc.h"
#include <vector>
#include "s2sClientIf.h"
#include "packet.h"

void marshalSubFilter(Pack &pk, const SubFilter &filter)
{
    pk << filter.interestedName
       << uint32_t(filter.interestedGroup)
       << uint32_t(filter.s2sType);
}

void unmarshalSubFilter(const Unpack &up, SubFilter &filter)
{
    std::string name;
    uint32_t group, s2sType;
    up >> name >> group >> s2sType;

    filter.interestedName = name;
    filter.interestedGroup = group;
    filter.s2sType = s2sType;
}

struct PSubFilter : public Marshallable
{
    std::vector<SubFilter> filters;

    virtual void marshal(Pack &pk) const
    {
        pk << uint32_t(filters.size());
        for(size_t i = 0; i < filters.size(); i++) {
            marshalSubFilter(pk, filters[i]);
        }
    }

    virtual void unmarshal(const Unpack &up)
    {
        uint32_t size = up.pop_uint32();
        filters.resize(size);
        for(size_t i = 0; i < size; i++) {
            unmarshalSubFilter(up, filters[i]);
        }
    }
};

void marshalS2SMeta(Pack &pk, const S2sMeta &meta)
{
    pk << uint64_t(meta.serverId) << uint32_t(meta.type) << meta.name << uint32_t(meta.groupId);
    pk.push_varstr32(meta.data);
    pk << uint64_t(meta.timestamp) << uint32_t(meta.status);
}

void unmarshalS2SMeta(const Unpack &up, S2sMeta &meta)
{
	uint64_t serverId;
	uint32_t type;
	std::string name;
	uint32_t groupId;
	std::string data;
	uint64_t timestamp;
	uint32_t status;

    up >> serverId >> type >> name >> groupId;
    data = up.pop_varstr32();
    up >> timestamp >> status;

    meta.serverId = serverId;
    meta.type = (MetaType)type;
    meta.name = name;
    meta.groupId = groupId;
    meta.data = data;
    meta.timestamp = timestamp;
    meta.status = (S2sMetaStatus)status;
}

struct PNotifyResult : public Marshallable
{
    uint32_t status; // S2sSessionStatus
    std::vector<S2sMeta> metas;

    virtual void marshal(Pack &pk) const
    {
        pk << status << uint32_t(metas.size());
        for(size_t i = 0; i < metas.size(); i++) {
            marshalS2SMeta(pk, metas[i]);
        }
    }

    virtual void unmarshal(const Unpack &up)
    {
        status = up.pop_uint32();
        uint32_t size = up.pop_uint32();
        metas.resize(size);
        for(size_t i = 0; i < size; i++) {
            unmarshalS2SMeta(up, metas[i]);
        }
    }
};

IMetaServer *s2sServer = NULL;

int initialize(const char *myName, const char *s2sKey, int myType)
{
    s2sServer = newMetaServer();
    if(s2sServer == NULL)
        return -1;
    return s2sServer->initialize(myName, s2sKey, MetaType(myType));
}

int subscribe(struct Buffer input)
{
    Unpack up(input.buffer, input.size);
    PSubFilter filter;
    filter.unmarshal(up);
    return s2sServer->subscribe(filter.filters);
}

struct Buffer pollNotify()
{
    PNotifyResult r;
    r.status = s2sServer->pollNotify(r.metas);

    PackBuffer pkbuf;
    Pack pk(pkbuf);
    r.marshal(pk);
    struct Buffer result;
    result.size = pkbuf.size();
    result.buffer = pkbuf.release();
    return result;
}

int setMine(const char *binData, int size)
{
    return s2sServer->setMine(std::string(binData, size));
}

int delMine()
{
    return s2sServer->delMine();
}

struct Buffer getMine()
{
    struct Buffer result;
    result.size = 0;
    result.buffer = NULL;

    S2sMeta meta;
    if(s2sServer->getMine(meta) == -1)
        return result;

    PackBuffer pkbuf;
    Pack pk(pkbuf);
    marshalS2SMeta(pk, meta);
    result.size = pkbuf.size();
    result.buffer = pkbuf.release();
    return result;
}
